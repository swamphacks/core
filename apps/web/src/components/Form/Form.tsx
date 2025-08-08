import { type FormProps, Form as RACForm } from "react-aria-components";
import { cn } from "@/utils/cn";
import { Suspense, type PropsWithChildren } from "react";
import { Spinner } from "@/components/ui/Spinner";

export function Form(props: FormProps & React.RefAttributes<HTMLFormElement>) {
  return (
    // We use suspense here since form fields are lazy loaded
    <FormSuspense>
      <RACForm {...props} className={cn("", props.className)} />
    </FormSuspense>
  );
}

export function FormSuspense({ children }: PropsWithChildren) {
  return (
    <Suspense
      fallback={
        <div className="flex w-full justify-center pt-10 gap-2 text-text-secondary">
          <Spinner />
          <p>Loading form...</p>
        </div>
      }
    >
      {children}
    </Suspense>
  );
}
