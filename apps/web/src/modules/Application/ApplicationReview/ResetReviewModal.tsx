import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { Group, Heading, Text } from "react-aria-components";

interface Props {
  isPending: boolean;
  onReset: () => Promise<void>;
}

const ResetReviewWarningModal = ({ isPending, onReset }: Props) => {
  return (
    <Modal isDismissible>
      <div className="flex flex-col gap-8 justify-center items-center">
        <div className="flex flex-col gap-4">
          <Heading className="text-lg" slot="title">
            Reset <strong>All</strong> Application Reviews?
          </Heading>
          <Text className="text-base">
            Are you sure you want to reset all application reviews? This action
            cannot be undone.
          </Text>
        </div>
        <div className="w-full flex flex-row justify-end">
          <Group className="gap-4 flex flex-row">
            <Button isPending={isPending} variant="secondary" slot="close">
              Nevermind
            </Button>

            <Button onPress={onReset} isPending={isPending} variant="danger">
              Reset
            </Button>
          </Group>
        </div>
      </div>
    </Modal>
  );
};

export default ResetReviewWarningModal;
