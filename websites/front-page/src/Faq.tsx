import { useState } from "react";
import "./FAQ.css";
import Camps from "./assets/camps.png";
import Fire from "./assets/Fire.gif";

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
        a: "Students 18 and older and from any university/college are eligible to participate. Students from all academic background are highly encouraged to participate. Regardless of your skill/experience level, we welcome all to join and make the most of the event!",
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
        q: "How long is SwampHacks?",
        a: "SwampHacks will be 36-hours beginning October 16 evening and going until October 18 afternoon.",
      },
      {
        q: "What I can I build?",
        a: "Anything from a chill video game for you and your friends to a supercharged sofa go kart! You have 36 hours to make something a reality, so make it count. It is important to note that your project does not have to be fully developed as you will be allowed to present *prototypes* for the product idea that you have developed.",
      },
      {
        q: "Will hardware be provided?",
        a: "Yes. We will have a hardware desk with a variety of devices you can borrow throughout the event. We will release more information about what we’ll have as the event approaches.",
      },
      {
        q: "How does project submission and judging work?",
        a: "To submit your project, we use Devpost and we will provide you the link as we get closer to the event. Judging will follow an expo format, where you will demo your project to the judges who visit your table. ",
      },
    ],
  },
  {
    section: "Logistics",
    faqs: [
      {
        q: "Is food provided?",
        a: "Yes, meals, snacks, and drinks will be provided throughout the hackathon to help keep you energized.",
      },
      {
        q: "Do I have to stay the whole time?",
        a: "You are welcome to leave and return to the venue(s) during the hackathon. However, we recommend to stay to maximize your time with your team, mentors, and the community. For overnight, we ask that you limit leaving and returning to the venue due to security.",
      },
      {
        q: "Are there showers available?",
        a: "There are not showers available inside of our venue(s). However, the J. Wayne Reitz Student Union has showers available on the Lower Level, which you can use during the open hours.",
      },
      {
        q: "Where can I sleep?",
        a: "If you would like to keep working on your project throughout the night, we will have a venue open to continue your collaboration and work. However, if you plan on sleeping through the night, we recommend that you head back to your place of residence.",
      },
    ],
  },
  {
    section: "Miscellaneous",
    faqs: [
      {
        q: "Looking to to sponsor?",
        a: "Please email our finance team, sponsors@swamphacks.com for more information. ",
      },
      {
        q: "Interested in volunteering?",
        a: "More details will come out as we get closer. See our instagram or discord server for all updates.",
      },
      {
        q: "Who can be a mentor?",
        a: "We welcome all who have prior hackathon experience, industry, or can assist those in need with basic project needs.",
      },
      {
        q: "How can I become a judge?",
        a: "We welcome individuals with prior hackathon experience, who are industry professionals or involved in academia. Please email contact@swamphacks.com with the subject “SwampHacks XII Judging Interest”.",
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
          <section className="faq-section" key={item.section}>
            <p className="faq-section-header">{item.section}</p>
            <div className="faq-section-list">
              {item.faqs.map((faqItem, questionIndex) => {
                const itemKey = `${sectionIndex}-${questionIndex}`;
                const isOpen = openQuestion === itemKey;

                return (
                  <div
                    className={`faq ${isOpen ? "open" : ""}`}
                    key={itemKey}
                    onClick={() => toggle(sectionIndex, questionIndex)}
                  >
                    <button
                      className={`faq-question ${isOpen ? "open" : ""}`}
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
          </section>
        ))}
      </div>

      <div className="campfire">
        <img className="camps" src={Camps} />
        <img className="fire" src={Fire} />
      </div>

      <a
        className="nes-btn is-primary register-button-2"
        href="https://app.swamphacks.com/application"
        target="_blank"
      >
        Apply Now!
      </a>
    </div>
  );
}
