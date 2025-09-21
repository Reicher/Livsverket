let hooks = [];
let currentHook = 0;
let rootComponent = null;
let rootContainer = null;

const SVG_TAGS = new Set(['svg', 'path']);

function scheduleRender() {
  if (!rootComponent || !rootContainer) {
    return;
  }
  currentHook = 0;
  const tree = rootComponent();
  rootContainer.innerHTML = '';
  rootContainer.appendChild(tree);
}

export function render(componentFn, container) {
  rootComponent = componentFn;
  rootContainer = container;
  hooks = [];
  scheduleRender();
}

export function useState(initialValue) {
  const hookIndex = currentHook;
  if (hooks.length <= hookIndex) {
    hooks.push(typeof initialValue === 'function' ? initialValue() : initialValue);
  }
  const setState = (value) => {
    const nextValue = typeof value === 'function' ? value(hooks[hookIndex]) : value;
    if (!Object.is(nextValue, hooks[hookIndex])) {
      hooks[hookIndex] = nextValue;
      scheduleRender();
    }
  };
  const value = hooks[hookIndex];
  currentHook += 1;
  return [value, setState];
}

export function createElement(tag, props = {}, children = []) {
  const namespace = SVG_TAGS.has(tag) ? 'http://www.w3.org/2000/svg' : undefined;
  const element = namespace ? document.createElementNS(namespace, tag) : document.createElement(tag);

  if (props.className) {
    element.setAttribute('class', props.className);
  }
  if (props.style) {
    Object.assign(element.style, props.style);
  }
  if (props.onClick) {
    element.addEventListener('click', props.onClick);
  }
  if (props.attrs) {
    Object.entries(props.attrs).forEach(([key, value]) => {
      element.setAttribute(key, value);
    });
  }
  const childArray = Array.isArray(children) ? children : [children];
  childArray.forEach((child) => {
    if (child == null) return;
    if (typeof child === 'string' || typeof child === 'number') {
      const node = namespace ? document.createTextNode(String(child)) : document.createTextNode(String(child));
      element.appendChild(node);
    } else {
      element.appendChild(child);
    }
  });
  return element;
}
