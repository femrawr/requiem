const tabs = document.querySelectorAll('.tab-button');
const contents = document.querySelectorAll('.content');

const createCheckbox = (config) => {
    const box = document.createElement('div');
    box.className = 'box';

    const title = document.createElement('div');
    title.className = 'title';

    const id = config.name.toLowerCase().replaceAll(' ', '_');

    const input = document.createElement('input');
    input.type = 'checkbox';
    input.id = config.id || id;
    input.checked = config.value ?? false;

    const label = document.createElement('label');
    label.innerHTML = `<strong>${config.name}</strong>`;

    title.appendChild(input);
    title.appendChild(label);

    const desc = document.createElement('div');
    desc.className = 'desc';
    desc.textContent = config.info;

    box.appendChild(title);
    box.appendChild(desc);

    return box;
};

const createTextbox = (config) => {
    const box = document.createElement('div');
    box.className = 'box';

    const title = document.createElement('div');
    title.className = 'title';

    const id = config.name.toLowerCase().replaceAll(' ', '_');

    const label = document.createElement('label');
    label.innerHTML = `<strong>${config.name}</strong>`;

    title.appendChild(label);

    const desc = document.createElement('div');
    desc.className = 'desc';
    desc.textContent = config.info;

    const container = document.createElement('div');
    container.className = 'input-container';

    const input = document.createElement('input');
    input.className = 'input';
    input.id = config.id || id;
    input.type = 'text';
    input.value = config.value ?? '';

    if (config.attributes) {
        for (const [key, val] of Object.entries(config.attributes)) {
            input.setAttribute(key, val);
        }
    }

    container.appendChild(input);

    box.appendChild(title);
    box.appendChild(desc);
    box.appendChild(container);

    return box;
};

const createDropdown = (config) => {
    const box = document.createElement('div');
    box.className = 'box';

    const title = document.createElement('div');
    title.className = 'title';

    const id = config.name.toLowerCase().replaceAll(' ', '_');

    const label = document.createElement('label');
    label.innerHTML = `<strong>${config.name}</strong>`;

    title.appendChild(label);

    const desc = document.createElement('div');
    desc.className = 'desc';
    desc.textContent = config.info;

    const container = document.createElement('div');
    container.className = 'input-container';

    const select = document.createElement('select');
    select.className = 'input';
    select.id = config.id || id;

    if (config.items && Array.isArray(config.items)) {
        config.items.forEach(item => {
            const option = document.createElement('option');
            option.value = item;
            option.textContent = item;
            select.appendChild(option);
        });
    }

    container.appendChild(select);

    box.appendChild(title);
    box.appendChild(desc);
    box.appendChild(container);

    return box;
};

const createSeparator = () => {
    const separator = document.createElement('div');
    separator.className = 'separator';

    return separator;
};

const createSetting = (config) => {
    switch (config.type) {
        case 'separator':
            return createSeparator();

        case 'checkbox':
            return createCheckbox(config);

        case 'text':
            return createTextbox(config);

        case 'list':
            return createDropdown(config);
    }
};

tabs.forEach((tab) => {
    tab.addEventListener('click', () => {
        tabs.forEach((button) => button.classList.remove('active'));
        contents.forEach((content) => content.classList.remove('active'));

        tab.classList.add('active');

        const target = tab.getAttribute('data-tab');
        document.getElementById(`${target}-tab`).classList.add('active');
    });
});

document.addEventListener('DOMContentLoaded', async () => {
    const fetchedLayout = await fetch('layout.json');
    const layout = await fetchedLayout.json();

    for (const [tab, item] of Object.entries(layout)) {
        const parent = document.getElementById(`${tab}-tab`);

        item.forEach((config) => {
            const box = createSetting(config);
            parent.appendChild(box);
        });
    }

    const fetchedCommands = await fetch('/api/get-commands');
    const commands = await fetchedCommands.json();

    const commandsParent = document.getElementById('commands-tab');

    commands.forEach((config) => {
        const box = createSetting({
            ...config,
            type: 'checkbox',
            value: true
        });

        commandsParent.appendChild(box);
    });

    document.querySelectorAll('input').forEach((input) => {
        input.setAttribute('spellcheck', 'false');
        input.setAttribute('autocomplete', 'off');
    });

    loadConfig();

    const inputs = document.querySelectorAll('input');
    const selects = document.querySelectorAll('select');

    [...inputs, ...selects].forEach((input) => {
        input.addEventListener('input', () => {
            localStorage.setItem(CFG_NAME, getConfig());
        });
    });
});
