import type { ComponentProps } from "react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type GradientButtonProps = ComponentProps<typeof Button>;

export function GradientButton({ className, ...props }: GradientButtonProps) {
  return (
    <Button
      className={cn(
        "bg-[color:var(--adfluence-orange)] text-white shadow-sm hover:bg-[color:var(--adfluence-orange-hover)]",
        className,
      )}
      {...props}
    />
  );
}
