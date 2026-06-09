interface Props {
  numAssigned: number;
  total: number;
}

export default function SummaryFooter({ numAssigned, total }: Props) {
  return (
    <div className="flex flex-col items-end text-text-secondary text-sm my-4">
      <p>Assigned Applications: {numAssigned}</p>
      <p>Total Applications: {total}</p>
      <p>Remaining: {total - numAssigned}</p>
    </div>
  );
}
