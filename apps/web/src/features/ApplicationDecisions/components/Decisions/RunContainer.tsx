// import RunCard from "./RunCard";

// interface Props {
//   eventId: string;
// }

// export default function RunContainer({
//   eventId,
// }: Props) {
//   const [runs, setRuns] = useState<Array<RunObject>>([]);

//   // const { addRun } = useRuns(eventId);

//   return (
//     <div>
//       <p className="mb-8 text-text-secondary">
//         Decision Runs:
//       </p>
//       <div className="flex flex-col gap-6 w-full">
//         {runs.length === 0 ? (
//           <p>No runs yet.</p>
//         ) : (
//           {runs.map((run) => {
//             <RunCard
//               key={run.id}
//               runData={run}
//             />
//           })}
//         )}
//       </div>
//     </div>
//   )
// }
