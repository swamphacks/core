import { useFormContext } from "@/components/Form/formContext";
import { Button } from "@/components/ui/Button";

export function SubmitButton({ label = "Submit" }: { label?: string }) {
  const form = useFormContext();
  return (
    <form.Subscribe
      selector={(state) => state.canSubmit}
      children={(canSubmit) => (
        <div className="flex justify-start mb-20">
          <Button type="submit" isDisabled={!canSubmit}>
            {label}
          </Button>
        </div>
      )}
    />
  );
}
