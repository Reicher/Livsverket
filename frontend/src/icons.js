import { createElement } from './mini-react.js';

function createMaterialIcon(name) {
  return createElement(
    'span',
    {
      className: 'material-symbols-rounded',
      attrs: { 'aria-hidden': 'true' }
    },
    name
  );
}

export function hierarchyIcon() {
  return createMaterialIcon('cruelty_free');
}

export function checklistIcon() {
  return createMaterialIcon('checklist');
}

export function settingsIcon() {
  return createMaterialIcon('settings');
}
