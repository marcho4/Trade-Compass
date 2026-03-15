import { cn } from "@/lib/utils";
import React from "react";

export default function GridBackground({
  children,
  className,
}: {
  children?: React.ReactNode;
  className?: string;
}) {
  return (
    <div className={cn("relative flex w-full items-center justify-center bg-background", className)}>
      <div
        className={cn(
          "absolute inset-0",
          "bg-size-[40px_40px]",
          "bg-[linear-gradient(to_right,oklch(0.9_0.005_240)_1px,transparent_1px),linear-gradient(to_bottom,oklch(0.9_0.005_240)_1px,transparent_1px)]",
          "dark:bg-[linear-gradient(to_right,oklch(0.3_0.015_240)_1px,transparent_1px),linear-gradient(to_bottom,oklch(0.3_0.015_240)_1px,transparent_1px)]",
        )}
      />
      <div className="pointer-events-none absolute inset-0 flex items-center justify-center bg-background mask-[radial-gradient(ellipse_at_center,transparent_20%,black)]"></div>
      {children}
    </div>
  );
}

