import type { ReactNode } from "react";

export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <div className="relative flex min-h-svh flex-col items-center justify-center overflow-hidden bg-[#06111F] p-4 sm:p-6">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_18%_24%,rgba(11,34,64,0.90)_0%,rgba(11,34,64,0.55)_34%,rgba(11,34,64,0.12)_58%,transparent_72%),radial-gradient(circle_at_82%_18%,rgba(124,58,237,0.20)_0%,rgba(124,58,237,0.10)_28%,transparent_58%),radial-gradient(circle_at_80%_76%,rgba(255,122,0,0.10)_0%,transparent_45%)]" />
      <div className="absolute inset-0 bg-[linear-gradient(92deg,rgba(6,17,31,0.98)_0%,rgba(7,26,47,0.96)_30%,rgba(11,34,64,0.84)_54%,rgba(11,34,64,0.40)_76%,rgba(11,34,64,0.08)_100%)]" />
      <div className="relative z-10 w-full max-w-2xl">{children}</div>
    </div>
  );
}
