import { createFileRoute } from "@tanstack/react-router";
import { Heading, DialogTrigger } from "react-aria-components";

import { Button } from "../../components/ui/Button";
import { AddEventModal } from "@/features/Admin/EventManager/components/AddEventModal";
import { useAdminEvents } from "@/features/Admin/EventManager/hooks/useAdminEvents";
import { EventDetailsCard } from "@/features/Admin/EventManager/components/EventDetailsCard";

export const Route = createFileRoute("/admin/events-management")({
  component: RouteComponent,
});

function RouteComponent() {
  const { data } = useAdminEvents();

  return (
    <div>
      <div className="mx-auto">
        <header className="flex justify-between items-center mb-6">
          <Heading className="text-text-main text-2xl">Event Manager</Heading>

          <DialogTrigger>
            <Button variant="primary">Create New Event</Button>

            <AddEventModal />
          </DialogTrigger>
        </header>
      </div>

      <div className="flex flex-row flex-wrap gap-6">
        {data &&
          data.map((val) => <EventDetailsCard key={val.id} event={val} />)}
      </div>
    </div>
  );

  // return (
  //   <div>
  //     <div className="max-w-1xl mx-auto">
  //       <header className="flex justify-between items-center mb-6">
  //         <h1 className="text-2xl font-bold">Events Management</h1>
  //         <Button variant="primary" onPress={() => setAddModalOpen(true)}>
  //           Add Event
  //         </Button>
  //       </header>

  //       {/* List of Event Cards */}
  //       <div className="space-y-4">
  //         {events.length > 0 ? (
  //           events.map((event) => (
  //             <AriaButton
  //               key={event.id}
  //               onPress={() => handleOpenEditModal(event)}
  //               className="w-full text-left focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 rounded-lg"
  //             >
  //               <Card className="p-6 w-full hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
  //                 <h2 className="font-bold text-lg">{event.name}</h2>
  //                 <p className="text-sm text-gray-600 dark:text-gray-400">
  //                   Starts: {new Date(event.startTime).toLocaleString()}
  //                 </p>
  //               </Card>
  //             </AriaButton>
  //           ))
  //         ) : (
  //           <Card className="p-6">
  //             <p>No events found. Click "Add Event" to create one.</p>
  //           </Card>
  //         )}
  //       </div>
  //     </div>

  //     {/* Modal for adding a new event (no changes) */}
  //     <ModalOverlay
  //       isOpen={isAddModalOpen}
  //       onOpenChange={setAddModalOpen}
  //       className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center"
  //     >
  //       <Modal className="w-full max-w-md outline-none">
  //         <Dialog className="outline-none">
  //           {({ close }) => (
  //             <Card className="p-6">
  //               <Heading slot="title" className="text-xl font-semibold mb-4">
  //                 Create New Event
  //               </Heading>
  //               <form className="space-y-4">
  //                 {/* Form inputs remain the same */}
  //                 <div>
  //                   <label
  //                     htmlFor="name"
  //                     className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1"
  //                   >
  //                     Event Name
  //                   </label>
  //                   <input
  //                     id="name"
  //                     name="name"
  //                     type="text"
  //                     value={eventData.name}
  //                     onChange={handleInputChange}
  //                     className="w-full p-2 rounded-md bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-600"
  //                   />
  //                 </div>
  //                 <div>
  //                   <label
  //                     htmlFor="applicationOpen"
  //                     className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1"
  //                   >
  //                     Application Open
  //                   </label>
  //                   <input
  //                     id="applicationOpen"
  //                     name="applicationOpen"
  //                     type="datetime-local"
  //                     value={eventData.applicationOpen}
  //                     onChange={handleInputChange}
  //                     className="w-full p-2 rounded-md bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-600"
  //                   />
  //                 </div>
  //                 <div>
  //                   <label
  //                     htmlFor="applicationClose"
  //                     className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1"
  //                   >
  //                     Application Close
  //                   </label>
  //                   <input
  //                     id="applicationClose"
  //                     name="applicationClose"
  //                     type="datetime-local"
  //                     value={eventData.applicationClose}
  //                     onChange={handleInputChange}
  //                     className="w-full p-2 rounded-md bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-600"
  //                   />
  //                 </div>
  //                 <div>
  //                   <label
  //                     htmlFor="startTime"
  //                     className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1"
  //                   >
  //                     Event Start Time
  //                   </label>
  //                   <input
  //                     id="startTime"
  //                     name="startTime"
  //                     type="datetime-local"
  //                     value={eventData.startTime}
  //                     onChange={handleInputChange}
  //                     className="w-full p-2 rounded-md bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-600"
  //                   />
  //                 </div>
  //                 <div>
  //                   <label
  //                     htmlFor="endTime"
  //                     className="block text-sm font-medium text-gray-500 dark:text-gray-400 mb-1"
  //                   >
  //                     Event End Time
  //                   </label>
  //                   <input
  //                     id="endTime"
  //                     name="endTime"
  //                     type="datetime-local"
  //                     value={eventData.endTime}
  //                     onChange={handleInputChange}
  //                     className="w-full p-2 rounded-md bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-600"
  //                   />
  //                 </div>
  //                 <div className="flex justify-end gap-2 pt-4">
  //                   <Button variant="secondary" onPress={close}>
  //                     Cancel
  //                   </Button>
  //                   <Button variant="primary" onPress={handleCreateEvent}>
  //                     Create Event
  //                   </Button>
  //                 </div>
  //               </form>
  //             </Card>
  //           )}
  //         </Dialog>
  //       </Modal>
  //     </ModalOverlay>

  //     {/* General Modal for Editing an Event */}
  //     <ModalOverlay
  //       isOpen={isEditModalOpen}
  //       onOpenChange={setEditModalOpen}
  //       className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center"
  //     >
  //       <Modal className="w-full max-w-md outline-none">
  //         <Dialog className="outline-none">
  //           {({ close }) => (
  //             <Card className="p-6">
  //               <Heading slot="title" className="text-xl font-semibold mb-2">
  //                 Edit Event
  //               </Heading>
  //               <p className="mb-6 text-lg text-gray-600 dark:text-gray-300">
  //                 {selectedEvent?.name}
  //               </p>

  //               <div className="space-y-3">
  //                 {/* Example Action Button */}
  //                 <Button
  //                   variant="primary"
  //                   onPress={handleAddUsers}
  //                   className="w-full"
  //                 >
  //                   Manage / Add Users
  //                 </Button>
  //                 {/* Other action buttons would go here */}
  //               </div>

  //               <div className="mt-6 pt-4 border-t border-gray-200 dark:border-gray-700 flex justify-between items-center">
  //                 {/* Delete button is now a secondary action */}
  //                 <Button variant="danger" onPress={handleDeleteEvent}>
  //                   Delete Event
  //                 </Button>
  //                 <Button variant="secondary" onPress={close}>
  //                   Close
  //                 </Button>
  //               </div>
  //             </Card>
  //           )}
  //         </Dialog>
  //       </Modal>
  //     </ModalOverlay>
  //   </div>
  // );
}
