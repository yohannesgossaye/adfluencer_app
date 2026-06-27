import type { ComponentProps } from "react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type GradientButtonProps = ComponentProps<typeof Button>;

export function GradientButton({ className, ...props }: GradientButtonProps) {
  return (
    <Button
      className={cn(
        "adfluence-gradient text-white shadow-md hover:opacity-95",
        className,
      )}
      {...props}
    />
  );
}
