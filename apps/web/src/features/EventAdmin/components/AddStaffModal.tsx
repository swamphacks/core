import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { TextField } from "@/components/ui/TextField";
import { useEventStaffUsers } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";
import { useUsers } from "@/features/Users/hooks/useUsers";
import { useContext, useMemo, useState } from "react";
import debounce from "lodash.debounce";
import { cn } from "@/utils/cn";
import type { User } from "@/lib/openapi/types";
import { useStaffActions } from "../hooks/useStaffActions";
import { toast } from "react-toastify";
import { OverlayTriggerStateContext } from "react-aria-components";

interface Props {
  eventId: string;
}

const AddStaffModal = ({ eventId }: Props) => {
  const state = useContext(OverlayTriggerStateContext)!;

  const [query, setQuery] = useState<string | null>("");
  const [selectedUsers, setSelectedUsers] = useState<User[]>([]);

  const {
    addMany: { mutate: addManyMutate, isPending: isAddPending },
  } = useStaffActions(eventId);

  const {
    data: users,
    isLoading: isUsersLoading,
    isError: isUsersError,
  } = useUsers(query);
  const {
    data: eventStaff,
    isLoading: isStaffUsersLoading,
    isError: isStaffUsersError,
  } = useEventStaffUsers(eventId);

  const handleSearch = useMemo(
    () =>
      debounce((value: string) => {
        setQuery(value);
      }, 500),
    [],
  );

  const onSelectionChange = (user: User) => {
    setSelectedUsers((prev) => {
      if (prev.some((u) => u.id === user.id)) {
        return prev.filter((u) => u.id !== user.id);
      } else {
        return [...prev, user];
      }
    });
  };

  const isAlreadyStaff = (userId: string) =>
    eventStaff?.some((staff) => staff.id === userId);

  const isSelected = (userId: string) =>
    selectedUsers.some((u) => u.id === userId);

  const displayedUsers: User[] = [
    ...selectedUsers,
    ...(users?.filter((u) => !selectedUsers.some((s) => s.id === u.id)) || []),
  ];

  const onAddStaff = async () => {
    if (selectedUsers.length === 0) return;

    addManyMutate(
      selectedUsers.map((u) => u.id),
      {
        onSuccess: () => {
          state.close();
          setSelectedUsers([]);
          setQuery("");
          toast.success("Staff members added successfully.");
        },
        onError: () => {
          toast.error("Failed to add staff members. Please try again.");
        },
      },
    );
  };

  return (
    <Modal size="2xl" title="Add Staff" className="h-1/2">
      <div className="w-full h-full justify-between py-4 flex flex-col gap-6">
        <div className="flex-1 flex flex-col gap-6">
          <TextField
            className="w-full"
            aria-label="Search users"
            placeholder="Search by name or email..."
            onChange={(e) => handleSearch(e)}
          />

          {isUsersLoading || isStaffUsersLoading ? (
            <p>Loading...</p>
          ) : isUsersError || isStaffUsersError ? (
            <p>Error loading users.</p>
          ) : displayedUsers.length === 0 ? (
            <p className="text-muted-foreground">No users found.</p>
          ) : (
            <div className="flex flex-col gap-4 max-h-96 overflow-y-auto">
              {displayedUsers.map((user) => {
                const alreadyStaff = isAlreadyStaff(user.id);
                const selected = isSelected(user.id);

                return (
                  <div
                    key={user.id}
                    onClick={() =>
                      !alreadyStaff && !isAddPending && onSelectionChange(user)
                    }
                    className={cn(
                      "w-full bg-input-bg hover:cursor-pointer hover:bg-neutral-900 rounded-md p-4 flex flex-row items-center justify-between select-none",
                      selected
                        ? "border-1 border-green-700"
                        : alreadyStaff
                          ? "opacity-50 cursor-not-allowed hover:cursor-default hover:bg-input-bg"
                          : "border-1 border-transparent",
                    )}
                  >
                    <div className="flex flex-col">
                      <p className="font-medium">{user.name}</p>
                      <p className="text-sm text-text-secondary">
                        {user.email}
                      </p>
                    </div>

                    {alreadyStaff && (
                      <p className="text-sm text-text-secondary font-medium">
                        Already Staff
                      </p>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>

        {selectedUsers.length > 0 && (
          <Button
            onClick={onAddStaff}
            variant="primary"
            isPending={isAddPending}
            isDisabled={isAddPending}
          >
            Add {selectedUsers.length} Staff Member
            {selectedUsers.length > 1 ? "s" : ""}
          </Button>
        )}
      </div>
    </Modal>
  );
};

export default AddStaffModal;
