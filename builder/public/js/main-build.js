const build = document.querySelector('#build');

const originTitle = document.title;

let timer = null;

const setTitle = (title) => {
    document.title = originTitle + ' - ' + title;

    if (timer) {
        clearTimeout(timer);
    }

    timer = setTimeout(() => {
        document.title = originTitle;
        timer = null;
    }, 6 * 60 * 1000);
};

build.addEventListener('click', async (e) => {
    const config = await fetch('/api/update-config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: getConfig()
    });

    if (!config.ok) {
        const text = await config.text();

        notif(text + '\n\n' + 'See console for full error.', 'Failed to update config', 'error', 10);
        console.warn(text);
        return;
    }

    if (e.ctrlKey) {
        notif('Successfully updated config', 'Builder', 'success', 10);
        return;
    }

    const initial = notif('Building...', 'Builder', 'info', 900);
    setTitle('Building');

    const build = await fetch('/api/start-build', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: getConfig()
    });

    if (!build.ok) {
        delNotif(initial);
        setTitle('Failed');

        const text = await build.text();

        notif(text, 'Failed to build', 'error', 30);
        console.warn(text);
        return;
    }

    new Audio('/assets/build-done.mp3')
        .play()
        .catch((err) => console.warn('failed to play sound -', err));

    delNotif(initial);
    notif('Successfully built.', 'Builder', 'success', 90);
    setTitle('Success');
});

document.addEventListener('keydown', (e) => {
    if (!e.ctrlKey) {
        return;
    }

    build.style.backgroundColor = 'var(--accent-dark)';
    build.style.borderColor = 'var(--accent-dark)';
    build.innerHTML = 'update config';
});

document.addEventListener('keyup', (e) => {
    if (e.ctrlKey) {
        return;
    }

    build.style.backgroundColor = 'var(--accent)';
    build.style.borderColor = 'var(--accent)';
    build.innerHTML = 'build';
})
