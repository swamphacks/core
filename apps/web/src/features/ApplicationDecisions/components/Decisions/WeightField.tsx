import type { Dispatch, SetStateAction } from "react";
import { Input, NumberField } from "react-aria-components";

interface Props {
  name: string;
  value: number;
  setValue: Dispatch<SetStateAction<number>>;
}

export default function WeightField({ name, value, setValue }: Props) {
  return (
    <div className="flex flex-row gap-2">
      <p className="text-text-secondary">{name}</p>
      <NumberField
        aria-label={`Weight for ${name ?? "no-name"}`}
        className="ml-auto"
        value={value}
        onChange={setValue}
      >
        <Input className="border-border border-1 disabled:border-0 w-12 h-6 text-sm px-1 disabled:bg-input-bg-disbaled dark:disabled:bg-neutral-800 disabled:cursor-not-allowed" />
      </NumberField>
    </div>
  );
}
