
"use client";

import * as DialogPrimitive from "@radix-ui/react-dialog";
import type { ComponentProps } from "react";
import clsx from "clsx";

type DialogContentProps = ComponentProps<
  typeof DialogPrimitive.Content
>;

export function DialogContent({
  children,
  className,
  ...props
}: DialogContentProps) {
  return (
    <DialogPrimitive.Portal>
      <DialogPrimitive.Overlay
        className="fixed inset-0 z-50 bg-black/50"
      />

      <DialogPrimitive.Content
        {...props}
        className={clsx(
          "fixed left-1/2 top-1/2 z-50",
          "w-full max-w-md",
          "-translate-x-1/2 -translate-y-1/2",
          "rounded-xl bg-white p-6 shadow-xl outline-none",
          className
        )}
      >
        {children}
      </DialogPrimitive.Content>
    </DialogPrimitive.Portal>
  );
}
