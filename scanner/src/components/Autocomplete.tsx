import { createSignal, Show, For, createEffect } from 'solid-js';

type Option = { label: string; value: string };

interface AutocompleteProps {
  options: Option[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}

const Autocomplete = (props: AutocompleteProps) => {
  const [input, setInput] = createSignal('');
  const [open, setOpen] = createSignal(false);

  // Sync input with value prop
  createEffect(() => {
    const match = props.options.find(opt => opt.value === props.value);
    setInput(match ? match.label : props.value || '');
  });

  const filtered = () =>
    props.options.filter(
      opt =>
        opt.label.toLowerCase().includes(input().toLowerCase()) ||
        opt.value.toLowerCase().includes(input().toLowerCase())
    );

  let inputRef: HTMLInputElement | undefined;

  const handleSelect = (opt: Option) => {
    setInput(opt.label);
    props.onChange(opt.value);
    setOpen(false);
  };

  const handleBlur = () => {
    setTimeout(() => setOpen(false), 100);
  };

  return (
    <div class="dropdown w-full">
      <input
        ref={inputRef}
        tabIndex={0}
        class="input input-bordered w-full"
        type="text"
        value={input()}
        placeholder={props.placeholder}
        onInput={e => {
          setInput(e.currentTarget.value);
          setOpen(true);
        }}
        onFocus={() => setOpen(true)}
        onBlur={handleBlur}
        autocomplete="off"
      />
      <Show when={open() && filtered().length > 0}>
        <ul tabIndex={0} class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-full max-h-48 overflow-auto z-10 border border-base-200 mt-1">
          <For each={filtered()}>
            {opt => (
              <li
                class="cursor-pointer"
                onMouseDown={e => {
                  e.preventDefault();
                  handleSelect(opt);
                }}
              >
                <a>{opt.label}</a>
              </li>
            )}
          </For>
        </ul>
      </Show>
    </div>
  );
};

export default Autocomplete;
