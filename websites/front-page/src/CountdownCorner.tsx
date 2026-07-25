import { useState, useEffect } from "react";
import "./CountdownCorner.css";

interface TimeLeft {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
  done: boolean;
}

interface CountdownCornerProps {
  targetDate?: Date;
  label?: string;
}

function getTimeLeft(targetDate: Date): TimeLeft {
  const diff = +targetDate - +new Date();
  if (diff <= 0) {
    return { days: 0, hours: 0, minutes: 0, seconds: 0, done: true };
  }
  return {
    days: Math.floor(diff / (1000 * 60 * 60 * 24)),
    hours: Math.floor((diff / (1000 * 60 * 60)) % 24),
    minutes: Math.floor((diff / 1000 / 60) % 60),
    seconds: Math.floor((diff / 1000) % 60),
    done: false,
  };
}

export default function CountdownCorner({
  targetDate = new Date(Date.now() + 1000 * 60 * 60 * 24 * 3),
  label = "Swamphacks Begins In:",
}: CountdownCornerProps) {
  const [timeLeft, setTimeLeft] = useState<TimeLeft>(() => getTimeLeft(targetDate));

  useEffect(() => {
    const interval = setInterval(() => {
      setTimeLeft(getTimeLeft(targetDate));
    }, 1000);
    return () => clearInterval(interval);
  }, [targetDate]);

  if (timeLeft.done) return null;

  const pad = (n: number) => String(n).padStart(2, "0");

  return (
    <div className="countdown-corner">
      <span className="countdown-corner__label">{label}</span>
      <div className="countdown-corner__digits">
        <span>{pad(timeLeft.days)}d</span>
        <span>{pad(timeLeft.hours)}h</span>
        <span>{pad(timeLeft.minutes)}m</span>
        <span>{pad(timeLeft.seconds)}s</span>
      </div>
    </div>
  );
}