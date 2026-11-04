"use client";

import { cn } from "@/lib/utils";

export type ChipOption<T extends string> = { value: T; label: string };

/**
 * One row of filter chips (a radio group): pick one value. Wraps on phones,
 * 44px tap targets on touch, labelled for screen readers.
 */
export function FilterChips<T extends string>({
  label,
  options,
  value,
  onChange,
  showLabel = false,
}: {
  label: string;
  options: ReadonlyArray<ChipOption<T>>;
  value: T;
  onChange: (value: T) => void;
  showLabel?: boolean;
}) {
  return (
    <div role="radiogroup" aria-label={label} className="flex flex-wrap items-center gap-2">
      {showLabel ? <span className="mr-1 text-xs font-medium text-[#94A3B8]">{label}</span> : null}
      {options.map((o) => {
        const active = o.value === value;
        return (
          <button
            key={o.value}
            type="button"
            role="radio"
            aria-checked={active}
            onClick={() => onChange(o.value)}
            className={cn(
              "inline-flex min-h-11 items-center rounded-full border px-4 text-sm font-medium transition-colors focus-visible:outline-2 focus-visible:outline-[#A78BFA] sm:min-h-9 sm:px-3.5 sm:text-[13px]",
              active
                ? "border-transparent bg-[color:var(--adfluence-action)] text-white"
                : "border-white/10 bg-white/5 text-[#CBD5E1] hover:border-white/20 hover:text-white",
            )}
          >
            {o.label}
          </button>
        );
      })}
    </div>
  );
}
