const getConfig = () => {
    const config = {
        tag: localStorage.getItem(TAG_NAME)
    };

    const inputs = document.querySelectorAll('input');
    const selects = document.querySelectorAll('select');

    [...inputs, ...selects].forEach((input) => {
        if (!input.id) {
            return;
        }

        if (input.type === 'checkbox') {
            config[input.id] = input.checked;
        } else {
            config[input.id] = input.value.trim();
        }
    });

    return JSON.stringify(config);
};

const loadConfig = () => {
    const saved = localStorage.getItem(CFG_NAME);
    if (!saved) {
        return;
    }

    const config = JSON.parse(saved);

    const inputs = document.querySelectorAll('input');
    const selects = document.querySelectorAll('select');

    [...inputs, ...selects].forEach((input) => {
        if (!input.id) {
            return;
        }

        if (!(input.id in config)) {
            return;
        }

        if (input.type === 'checkbox') {
            input.checked = !!config[input.id];
        } else {
            input.value = config[input.id] || '';
        }
    });
};