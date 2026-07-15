import { useState } from "react";
import "./FAQ.css";

const faqSections = [
  {
    section: "General",
    faqs: [
      {
        q: "What is SwampHacks?",
        a: "SwampHacks is a 36-hour event, where participants collaborate in teams to create solutions, typically in the form of software and hardware projects. Participants can work in teams or individually to develop and build innovative solutions.",
      },
      {
        q: "Who can participate?",
        a: "Students 18 and older and from any university/college are eligible to participate. Students from all academic background are highly encouraged to participate. Regardless of your skill/experience level, we welcome all to join and make the most of the event! As part of our dedication to the UF community and fostering experiential learning in early computing education, we’ll be releasing more details soon about our application review and acceptance processes.",
      },
      {
        q: "Do I need to have coding experience to participate?",
        a: "No experience, no problem! We welcome participant of all skill levels. The idea is that throughout the event, we will support you through workshops and mentors available around the clock.",
      },
      {
        q: "It’s my first hackathon. Should I apply?",
        a: "Yes! We highly encourage and want first time hackers to apply. We’re dedicated to providing experiential learning opportunities for those early in their computing journey.",
      },
      {
        q: "How can I prepare for the hackathon?",
        a: "You can prepare by familiarizing yourself with coding languages or tools, you are looking to use/explore, brushing up on project management skills, and collaboration. However, everything you need will be available during the hackathon.",
      },
      {
        q: "How much does it cost to participate?",
        a: "Nothing! SwampHacks is completely free for all student accepted to the event. However, SwampHacks is unable to cover transportation costs at this time.",
      },
      {
        q: "Is SwampHacks in-person or virtual?",
        a: "SwampHacks is FULLY in-person! Everything from check-in to judging will occur at our venue. Make sure you are able to attend in-person before registering.",
      },
    ],
  },
  {
    section: "Event Details",
    faqs: [
      {
        q: "What is SwampHacks?",
        a: "SwampHacks is a 36-hour event, where participants collaborate in teams to create solutions, typically in the form of software and hardware projects. Participants can work in teams or individually to develop and build innovative solutions.",
      },
      {
        q: "Who can participate?",
        a: "Students 18 and older and from any university/college are eligible to participate. Students from all academic background are highly encouraged to participate. Regardless of your skill/experience level, we welcome all to join and make the most of the event! As part of our dedication to the UF community and fostering experiential learning in early computing education, we’ll be releasing more details soon about our application review and acceptance processes.",
      },
      {
        q: "Do I need to have coding experience to participate?",
        a: "No experience, no problem! We welcome participant of all skill levels. The idea is that throughout the event, we will support you through workshops and mentors available around the clock.",
      },
      {
        q: "It’s my first hackathon. Should I apply?",
        a: "Yes! We highly encourage and want first time hackers to apply. We’re dedicated to providing experiential learning opportunities for those early in their computing journey.",
      },
      {
        q: "How can I prepare for the hackathon?",
        a: "You can prepare by familiarizing yourself with coding languages or tools, you are looking to use/explore, brushing up on project management skills, and collaboration. However, everything you need will be available during the hackathon.",
      },
      {
        q: "How much does it cost to participate?",
        a: "Nothing! SwampHacks is completely free for all student accepted to the event. However, SwampHacks is unable to cover transportation costs at this time.",
      },
      {
        q: "Is SwampHacks in-person or virtual?",
        a: "SwampHacks is FULLY in-person! Everything from check-in to judging will occur at our venue. Make sure you are able to attend in-person before registering.",
      },
    ],
  },
];

export default function FAQ() {
  const [openQuestion, setOpenQuestion] = useState<string | null>(null);

  const toggle = (sectionIndex: number, questionIndex: number) => {
    const key = `${sectionIndex}-${questionIndex}`;
    setOpenQuestion(openQuestion === key ? null : key);
  };

  return (
    <div className="faq-container">
      <p className="faq-header">FAQ</p>
      <div className="faq-items-container">
        {faqSections.map((item, sectionIndex) => (
          <div key={item.section}>
            <p>{item.section}</p>
            {item.faqs.map((faqItem, questionIndex) => {
              const itemKey = `${sectionIndex}-${questionIndex}`;
              const isOpen = openQuestion === itemKey;

              return (
                <div className={`faq ${isOpen ? "open" : ""}`} key={itemKey}>
                  <button
                    className={`faq-question ${isOpen ? "open" : ""}`}
                    onClick={() => toggle(sectionIndex, questionIndex)}
                    aria-expanded={isOpen}
                  >
                    {faqItem.q}
                  </button>

                  {isOpen && (
                    <div className="faq-answer">
                      <p>{faqItem.a}</p>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        ))}
      </div>
    </div>
  );
}
