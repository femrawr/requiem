const KEY_SEPARATOR = '!';

const getConfig = async (password) => {
    const config = {
        tag: localStorage.getItem(TAG_NAME)
    };

    const inputs = document.querySelectorAll('input');
    const selects = document.querySelectorAll('select');

    [...inputs, ...selects].forEach((input) => {
        if (input.type === 'checkbox') {
            config[input.id] = input.checked;
        } else {
            config[input.id] = input.value.trim();
        }
    });

    if (password) {
        const encryptedConfig = Object.entries(config).map(async ([name, value]) => {
            const newValue = await encryptData(JSON.stringify(value), password);
            return [name, newValue];
        });

        const newConfig = Object.fromEntries(
            await Promise.all(encryptedConfig)
        );

        return JSON.stringify(newConfig);
    }

    return JSON.stringify(config);
};

const loadConfig = async (password) => {
    const saved = localStorage.getItem(CFG_NAME);
    if (!saved) return;

    let config = JSON.parse(saved);

    if (password) {
        const decryptedConfig = Object.entries(config).map(async ([name, value]) => {
            const decrypted = await decryptData(String(value), password);

            try {
                return [name, JSON.parse(decrypted)];
            } catch {
                return [name, decrypted];
            }
        });

        config = Object.fromEntries(
            await Promise.all(decryptedConfig)
        );
    }

    const inputs = document.querySelectorAll('input');
    const selects = document.querySelectorAll('select');

    [...inputs, ...selects].forEach((input) => {
        if (input.type === 'checkbox') {
            input.checked = config[input.id] === true;
        } else {
            input.value = config[input.id] ?? '';
        }
    });
};
