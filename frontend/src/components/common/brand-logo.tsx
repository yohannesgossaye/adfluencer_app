import Image from "next/image";
import Link from "next/link";

import { cn } from "@/lib/utils";

type BrandLogoProps = {
  href?: string;
  className?: string;
  iconClassName?: string;
  showName?: boolean;
  /** Hide wordmark below `sm` (640px) — icon only on small phones */
  showNameFrom?: "always" | "sm";
  tagline?: string;
  iconSize?: "sm" | "md";
};

export function BrandLogo({
  href = "/",
  className,
  iconClassName,
  showName = true,
  showNameFrom = "always",
  tagline = "Ethiopia's creator marketplace",
  iconSize = "md",
}: BrandLogoProps) {
  const iconDimensions = iconSize === "sm" ? 32 : 36;

  const content = (
    <div className={cn("flex items-center gap-2.5", className)}>
      <Image
        src="/brand/adfluence_icon_nav.svg"
        alt=""
        width={iconDimensions}
        height={iconDimensions}
        className={cn(
          "block shrink-0 rounded-[10px]",
          iconSize === "sm" ? "size-8" : "size-9",
          iconClassName,
        )}
        priority
        aria-hidden
      />
      {showName ? (
        <div
          className={cn(
            "leading-tight",
            showNameFrom === "sm" && "hidden sm:block",
          )}
        >
          <p className="font-[family-name:var(--font-display)] text-xl font-extrabold tracking-tight text-white">Adfluence</p>
          <p className="text-xs text-[#CBD5E1]">{tagline}</p>
        </div>
      ) : null}
    </div>
  );

  if (!href) {
    return content;
  }

  return (
    <Link href={href} className="flex shrink-0 items-center">
      {content}
    </Link>
  );
}
