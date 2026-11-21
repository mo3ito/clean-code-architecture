
"use client";

import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger} from "@/src/shared/components/molecules/dialog"

// import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <main className="flex min-h-screen items-center justify-center">
      <Dialog>
        <DialogTrigger asChild>
          <button>Open Dialog</button>
        </DialogTrigger>

        <DialogContent>
          <DialogHeader>
            <DialogTitle>Welcome to Radix Dialog</DialogTitle>

            <DialogDescription>
              This is a reusable Dialog component built with Radix UI.
              You can use it for forms, confirmations, and other interactions.
            </DialogDescription>
          </DialogHeader>

          <DialogFooter>
            <DialogClose asChild>
              <button >Cancel</button>
            </DialogClose>

            <DialogClose asChild>
              <button>Confirm</button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </main>
  );
}
