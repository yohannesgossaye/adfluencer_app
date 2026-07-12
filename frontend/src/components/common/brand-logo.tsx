import Image from "next/image";
import Link from "next/link";

import { cn } from "@/lib/utils";

type BrandLogoProps = {
  variant?: "icon" | "full";
  href?: string;
  className?: string;
  iconClassName?: string;
  showName?: boolean;
  tagline?: string;
  iconSize?: "sm" | "md";
};

export function BrandLogo({
  variant = "icon",
  href = "/",
  className,
  iconClassName,
  showName = true,
  tagline = "Connect. Collaborate. Grow.",
  iconSize = "md",
}: BrandLogoProps) {
  const iconDimensions = iconSize === "sm" ? 32 : 36;

  const content =
    variant === "full" ? (
      <Image
        src="/brand/adfluence_primary_dark.svg"
        alt="Adfluence"
        width={160}
        height={40}
        className={cn("h-9 w-auto", className)}
        priority
      />
    ) : (
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
          <div className="leading-tight">
            <p className="text-xl font-bold text-white">Adfluence</p>
            <p className="text-xs text-slate-300">{tagline}</p>
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
