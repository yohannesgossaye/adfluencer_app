import { cn } from "@/lib/utils";

/**
 * Our signature element: a thin band inspired by the woven tibeb borders on
 * Ethiopian cloth. Used as a divider, never as a full background.
 */
export function TibebBand({ className }: { className?: string }) {
  return (
    <div aria-hidden className={cn("h-7 w-full overflow-hidden", className)}>
      <svg width="100%" height="28" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none">
        <defs>
          <pattern id="tibeb" width="144" height="28" patternUnits="userSpaceOnUse">
            <rect width="144" height="28" fill="#F3E7CF" />
            <rect width="144" height="3" fill="#4C1D95" />
            <rect y="25" width="144" height="3" fill="#4C1D95" />
            {[
              ["#4C1D95", "#B45309"],
              ["#B45309", "#2F7D4F"],
              ["#2F7D4F", "#4C1D95"],
            ].map(([main, small], i) => {
              const x = i * 48;
              const cx = x + 24;
              return (
                <g key={i}>
                  <path d={`M${cx} 5 L${cx + 9} 14 L${cx} 23 L${cx - 9} 14 Z`} fill={main} />
                  <path d={`M${cx} 10 L${cx + 4} 14 L${cx} 18 L${cx - 4} 14 Z`} fill="#F3E7CF" />
                  <path d={`M${x} 14 L${x + 6} 9 L${x + 12} 14 L${x + 6} 19 Z`} fill={small} opacity="0.85" />
                </g>
              );
            })}
          </pattern>
        </defs>
        <rect width="100%" height="28" fill="url(#tibeb)" />
      </svg>
    </div>
  );
}
