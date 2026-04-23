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
    let config = await getConfig();
    let parsedConfig = JSON.parse(config);

    if (parsedConfig.end_to_end_encryption === true) {
        const clientKeys = await generateKeyPair();

        const fetchedServerPublicKey = await fetch('/api/get-key', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ key: await exportPublicKey(clientKeys.publicKey) })
        });

        if (!fetchedServerPublicKey.ok) {
            const text = await fetchedServerPublicKey.text();

            notif(text + '\n\n' + 'See console for full error.', 'Failed to get server key.', NOTIF_ERROR, 10);
            console.warn(text);
            return;
        }

        const serverPublicKey = await fetchedServerPublicKey.text();
        const sharedSecret = await getSharedSecret(clientKeys.privateKey, serverPublicKey);

        config = await getConfig(sharedSecret);
    }

    const update = await fetch('/api/update-config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: config
    });

    if (!update.ok) {
        const text = await update.text();

        notif(text + '\n\n' + 'See console for full error.', 'Failed to update config', NOTIF_ERROR, 10);
        console.warn(text);
        return;
    }

    if (e.ctrlKey) {
        notif('Successfully updated config', 'Builder', NOTIF_SUCC, 10);
        return;
    }

    const initial = notif('Building...', 'Builder', NOTIF_INFO, 900);
    setTitle('Building');

    const build = await fetch('/api/start-build', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: config
    });

    if (!build.ok) {
        delNotif(initial);
        setTitle('Failed');

        const text = await build.text();

        notif(text, 'Failed to build', NOTIF_ERROR, 30);
        console.warn(text);
        return;
    }

    new Audio('/assets/build-done.mp3')
        .play()
        .catch((err) => console.warn('failed to play sound -', err));

    delNotif(initial);
    notif('Successfully built.', 'Builder', NOTIF_SUCC, 90);
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
