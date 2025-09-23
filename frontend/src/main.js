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

const HIERARCHY_DATA = {
  id: 'life',
  label: 'Life',
  children: [
    { id: 'bacteria', label: 'Bacteria', children: [] },
    { id: 'archaea', label: 'Archaea', children: [] },
    { id: 'eukarya', label: 'Eukarya', children: [] }
  ]
};

const BUTTONS = [
  { key: 'hierarchy', label: 'Hierarchy', icon: hierarchyIcon },
  { key: 'checklist', label: 'Checklist', icon: checklistIcon },
  { key: 'settings', label: 'Settings', icon: settingsIcon }
];

function createTopBar() {
  return createElement('header', { className: 'top-bar' }, 'Livsverket');
}

function createHierarchyNode(node, depth, expandedNodes, onToggleNode) {
  const hasChildren = node.children.length > 0;
  const isExpanded = expandedNodes[node.id] ?? (depth === 0);
  const toggleClasses = ['node-toggle'];
  if (!hasChildren) {
    toggleClasses.push('node-toggle--inactive');
  }

  const toggleButton = createElement(
    'button',
    {
      className: toggleClasses.join(' '),
      onClick: (event) => {
        event.stopPropagation();
        if (!hasChildren) {
          return;
        }
        onToggleNode(node.id);
      },
      attrs: {
        'aria-label': `${isExpanded ? 'Collapse' : 'Expand'} ${node.label}`,
        'aria-expanded': hasChildren ? (isExpanded ? 'true' : 'false') : 'false',
        'aria-disabled': hasChildren ? 'false' : 'true'
      }
    },
    hasChildren && isExpanded ? '−' : '+'
  );

  const header = createElement('div', { className: 'hierarchy-node-header' }, [
    toggleButton,
    createElement('span', { className: 'hierarchy-node-label' }, node.label)
  ]);

  const children =
    isExpanded && node.children.length > 0
      ? createElement(
          'div',
          { className: 'hierarchy-children' },
          node.children.map((child) =>
            createHierarchyNode(child, depth + 1, expandedNodes, onToggleNode)
          )
        )
      : null;

  return createElement(
    'div',
    {
      className: 'hierarchy-node',
      style: { marginLeft: `${depth * 12}px` }
    },
    [header, children]
  );
}

function createHierarchyContent(statusMessage, expandedNodes, onToggleNode) {
  const content = createElement('main', {
    className: 'content-area hierarchy-content'
  });
  content.appendChild(
    createElement('div', { className: 'hierarchy-tree' }, [
      createHierarchyNode(HIERARCHY_DATA, 0, expandedNodes, onToggleNode)
    ])
  );
  if (statusMessage) {
    content.appendChild(
      createElement('p', { className: 'status-message hierarchy-status' }, statusMessage)
    );
  }
  return content;
}

function createContent(activePage, statusMessage, expandedNodes, onToggleNode) {
  if (activePage === 'hierarchy') {
    return createHierarchyContent(statusMessage, expandedNodes, onToggleNode);
  }

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
  const buttonElement = createElement(
    'button',
    {
      className: 'nav-button',
      onClick: () => onSelect(button.key),
      attrs: { 'aria-label': button.label }
    },
    icon
  );
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
  const [expandedNodes, setExpandedNodes] = useState(() => ({ life: true }));

  if (status === 'Connecting to backend...') {
    fetchStatus().then(setStatus);
  }

  const toggleNode = (nodeId) => {
    setExpandedNodes((prev) => ({
      ...prev,
      [nodeId]: !prev[nodeId]
    }));
  };

  const shell = createElement('div', { className: 'app-shell' });
  shell.appendChild(createTopBar());
  shell.appendChild(createContent(activePage, status, expandedNodes, toggleNode));
  shell.appendChild(createBottomNav(activePage, setActivePage));
  return shell;
}

const container = document.getElementById('root');
render(App, container);
