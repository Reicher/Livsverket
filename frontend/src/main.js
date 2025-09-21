import { render, useState, createElement } from './mini-react.js';
import { hierarchyIcon, checklistIcon, settingsIcon } from './icons.js';

const PAGES = {
  hierarchy: {
    title: 'Hierarchy',
    description: 'Hierarchy content'
  },
  checklist: {
    title: 'Checklist',
    description: 'Checklist content'
  },
  settings: {
    title: 'Settings',
    description: 'Settings content'
  }
};

const BUTTONS = [
  { key: 'hierarchy', label: 'Hierarchy', icon: hierarchyIcon },
  { key: 'checklist', label: 'Checklist', icon: checklistIcon },
  { key: 'settings', label: 'Settings', icon: settingsIcon }
];

function createTopBar() {
  return createElement('header', { className: 'top-bar' }, 'Livsverket');
}

function createContent(activePage, statusMessage) {
  const { title, description } = PAGES[activePage];
  const content = createElement('main', { className: 'content-area' });
  const card = createElement('section', { className: 'content-card' }, [
    createElement('h1', {}, title),
    createElement('p', {}, description),
    statusMessage
      ? createElement('p', { className: 'status-message' }, statusMessage)
      : null
  ]);
  content.appendChild(card);
  return content;
}

function createNavButton(button, isActive, onSelect) {
  const icon = button.icon();
  icon.classList.add('nav-icon');
  const label = createElement('span', {}, button.label);
  const buttonElement = createElement('button', { className: 'nav-button', onClick: () => onSelect(button.key) }, [
    icon,
    label
  ]);
  if (isActive) {
    buttonElement.classList.add('active');
  }
  return buttonElement;
}

function createBottomNav(activePage, onSelect) {
  const nav = createElement('nav', { className: 'bottom-nav' });
  BUTTONS.forEach((button) => {
    nav.appendChild(createNavButton(button, button.key === activePage, onSelect));
  });
  return nav;
}

async function fetchStatus() {
  try {
    const response = await fetch('http://localhost:8080/api/status');
    if (!response.ok) {
      throw new Error('Request failed');
    }
    const data = await response.json();
    return data.message;
  } catch (error) {
    console.warn('Unable to reach backend:', error.message);
    return 'Backend connection unavailable';
  }
}

function App() {
  const [activePage, setActivePage] = useState('hierarchy');
  const [status, setStatus] = useState('Connecting to backend...');

  if (status === 'Connecting to backend...') {
    fetchStatus().then(setStatus);
  }

  const shell = createElement('div', { className: 'app-shell' });
  shell.appendChild(createTopBar());
  shell.appendChild(createContent(activePage, status));
  shell.appendChild(createBottomNav(activePage, setActivePage));
  return shell;
}

const container = document.getElementById('root');
render(App, container);
