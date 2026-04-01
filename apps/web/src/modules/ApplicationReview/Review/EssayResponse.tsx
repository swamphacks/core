interface EssayResponseProps {
  question: string;
  response: string;
}

export default function EssayResponse({
  question,
  response,
}: EssayResponseProps) {
  return (
    <div>
      <p className="text-text-secondary">{question}</p>
      <p>{response}</p>
    </div>
  );
}
