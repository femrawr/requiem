const tag = document.querySelector('.build-tag');
const saved = localStorage.getItem(TAG_NAME) || '';

tag.addEventListener('dblclick', (e) => {
    if (!e.ctrlKey) {
        return;
    }

    tag.contentEditable = true;
    tag.focus();
});

tag.addEventListener('blur', () => {
    tag.contentEditable = false;
    const text = tag.textContent.trim();

    if (text) {
        localStorage.setItem(TAG_NAME, text);
    } else {
        tag.textContent = '';
        localStorage.removeItem(TAG_NAME);
    }
});

tag.addEventListener('keydown', (e) => {
    if (e.key !== 'Enter') {
        return;
    }

    e.preventDefault();
    tag.blur();
});

document.addEventListener('DOMContentLoaded', () => {
    if (saved) {
        tag.textContent = saved;
    }

    if (tag.textContent !== 'default') {
        return;
    }

    tag.style.visibility = 'hidden';

    document.addEventListener('keydown', (e) => {
        if (e.ctrlKey) {
            tag.style.visibility = 'visible';
        }
    });
});
