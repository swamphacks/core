import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";
import { Button as RACButton } from "react-aria-components";
import { useToggleState } from "react-stately";
import TablerChevronDown from "~icons/tabler/chevron-down";
import TablerChevronUp from "~icons/tabler/chevron-up";

export const Route = createFileRoute("/_protected/_user/resources/programming")(
  {
    component: RouteComponent,
  },
);

type Boilerplate = {
  name: string;
  nameLink: string;
  technology: string;
  technologyLink: string;
  use: string;
  language: string;
};

type AccordionLink = {
  text: string;
  url: string;
};

type AccordionItem = {
  title: string;
  description?: string;
  links: AccordionLink[];
};

const boilerplateData: Boilerplate[] = [
  {
    name: "MLH Node.js Starter",
    nameLink: "https://github.com/MLH/mlh-hackathon-nodejs-starter",
    technology: "Node.js",
    technologyLink: "https://nodejs.org/en/docs/",
    use: "backend/web dev",
    language: "JavaScript",
  },
  {
    name: "Next.js Official Starter",
    nameLink: "https://nextjs.org/docs/app/getting-started/installation",
    technology: "Next.js",
    technologyLink: "https://nextjs.org/docs/app/getting-started",
    use: "web dev",
    language: "Javascript, TypeScript",
  },
  {
    name: "React Native Boilerplate by Infinite Red",
    nameLink: "https://github.com/infinitered/ignite",
    technology: "React Native",
    technologyLink: "https://reactnative.dev/",
    use: "mobile app dev",
    language: "TypeScript",
  },
  {
    name: "Arduino-ESP32 Getting Started Guide",
    nameLink:
      "https://docs.espressif.com/projects/arduino-esp32/en/latest/getting_started.html",
    technology: "Arduino-ESP32",
    technologyLink:
      "https://github.com/espressif/arduino-esp32?tab=readme-ov-file",
    use: "hardware dev, IoT",
    language: "C++",
  },
];

const coreFeaturesData: AccordionItem[] = [
  {
    title: "APIs and backend services",
    description: "Resources for building APIs and backend services:",
    links: [
      {
        text: "Express.js - Fast, unopinionated web framework",
        url: "https://expressjs.com/en/starter/installing.html",
      },
      {
        text: "FastAPI - Modern Python web framework",
        url: "https://fastapi.tiangolo.com/",
      },
      {
        text: "REST API Tutorial - Postman Learning Center",
        url: "https://www.postman.com/learn/rest-apis",
      },
    ],
  },
  {
    title: "Databases and data storage",
    description: "Learn about databases and data storage:",
    links: [
      {
        text: "MongoDB Getting Started",
        url: "https://www.mongodb.com/docs/manual/tutorial/getting-started/",
      },
      {
        text: "PostgreSQL Tutorial",
        url: "https://www.postgresql.org/docs/current/tutorial.html",
      },
      {
        text: "Firebase Firestore - NoSQL database",
        url: "https://firebase.google.com/docs/firestore",
      },
      {
        text: "Convex - Backend that is always in sync",
        url: "https://docs.convex.dev/home",
      },
      {
        text: "Node.js + Express + MongoDB Tutorials",
        url: "https://www.youtube.com/playlist?list=PLmIhXONtG4Wsj3SfwWvZLGjJzEYQzXdO8",
      },
      {
        text: "FreeCodeCamp - Intro to Backend Web Development",
        url: "https://www.youtube.com/watch?v=KOutPbKc9UM",
      },
    ],
  },
  {
    title: "Authentication and user management",
    description: "Authentication and user management resources:",
    links: [
      {
        text: "Auth0 - Authentication platform",
        url: "https://auth0.com/docs/get-started",
      },
      {
        text: "Firebase Authentication",
        url: "https://firebase.google.com/docs/auth",
      },
      {
        text: "NextAuth.js - Next.js authentication",
        url: "https://next-auth.js.org/getting-started/introduction",
      },
      {
        text: "Clerk - Complete web user management",
        url: "https://clerk.com/docs",
      },
    ],
  },
  {
    title: "Working with sensors, motors, or inputs",
    description: "Hardware integration guides:",
    links: [
      {
        text: "Adafruit Learning System - Free IoT project guides and ideas",
        url: "https://learn.adafruit.com/",
      },
      {
        text: "Basic breadboarding guide",
        url: "https://learn.sparkfun.com/tutorials/how-to-use-a-breadboard/all",
      },
      {
        text: "Arduino Tutorials",
        url: "https://www.arduino.cc/en/Tutorial/HomePage",
      },
      {
        text: "Raspberry Pi Getting Started",
        url: "https://www.raspberrypi.com/documentation/computers/getting-started.html",
      },
      {
        text: "ESP32/ESP8266 IoT Projects",
        url: "https://www.youtube.com/playlist?list=PLNFq0T6Z3JPsMWtVgmqPnpIu41LNMDPiA",
      },
    ],
  },
];

const debuggingData: AccordionItem[] = [
  {
    title: "Discord Support System",
    description: "Get help from our mentors and other hackers in real-time:",
    links: [
      {
        text: "Join the SwampHacks Discord server and create request in the #support channel",
        url: "https://discord.gg/qMKXgKhU",
      },
    ],
  },
  {
    title: "Debugging web and backend issues",
    description: "Tools and techniques for debugging web applications:",
    links: [
      {
        text: "Chrome DevTools Documentation",
        url: "https://developer.chrome.com/docs/devtools/",
      },
      {
        text: "React DevTools - Debug React applications",
        url: "https://react.dev/learn/react-developer-tools",
      },
      {
        text: "Node.js Debugging Guide",
        url: "https://nodejs.org/en/docs/guides/debugging-getting-started/",
      },
      {
        text: "Python Debugging with pdb",
        url: "https://docs.python.org/3/library/pdb.html",
      },
      {
        text: "VS Code Debugging Guide",
        url: "https://code.visualstudio.com/docs/editor/debugging",
      },
    ],
  },
  {
    title: "Hardware troubleshooting checklists",
    description: "Common hardware issues and how to fix them:",
    links: [
      {
        text: "Arduino Troubleshooting Guide",
        url: "https://www.arduino.cc/en/Guide/Troubleshooting",
      },
      {
        text: "Raspberry Pi Troubleshooting",
        url: "https://www.raspberrypi.com/documentation/computers/raspberry-pi.html#troubleshooting",
      },
      {
        text: "ESP32 Common Issues and Solutions",
        url: "https://randomnerdtutorials.com/esp32-troubleshooting-guide/",
      },
    ],
  },
  {
    title: "Common configuration mistakes",
    description: "Avoid these frequent setup and configuration errors:",
    links: [
      {
        text: "CORS Errors - How to Fix",
        url: "https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS/Errors",
      },
      {
        text: "Environment Variables Best Practices",
        url: "https://dev.to/khalidk799/environment-variables-its-best-practices-1o1o",
      },
      {
        text: "Port Already in Use - Solutions",
        url: "https://dev.to/antonrosh/address-already-in-use-a-simple-guide-to-freeing-up-ports-54g5",
      },
    ],
  },
  {
    title: "Logs, errors, and testing tools",
    description: "Tools to help you understand what's going wrong:",
    links: [
      {
        text: "Error Tracking with Sentry",
        url: "https://docs.sentry.io/platforms/javascript/",
      },
      {
        text: "Postman - API Testing and Debugging",
        url: "https://learning.postman.com/",
      },
      {
        text: "Using the debugger in VS Code",
        url: "https://code.visualstudio.com/docs/debugtest/debugging",
      },
      {
        text: "Debugging on the web with Chrome DevTools",
        url: "https://developer.chrome.com/docs/devtools/javascript",
      },
    ],
  },
];

function AccordionItem({
  title,
  description,
  links,
}: {
  title: string;
  description?: string;
  links: AccordionLink[];
}) {
  const toggleState = useToggleState({ defaultSelected: false });
  const isExpanded = toggleState.isSelected;

  return (
    <div className="border border-neutral-200 dark:border-neutral-800 rounded-md overflow-hidden mb-2">
      <RACButton
        onPress={toggleState.toggle}
        className="w-full px-4 py-3 flex items-center justify-between bg-surface"
      >
        <span className="font-semibold text-left">{title}</span>
        <span className="flex items-center justify-center transition-transform duration-150">
          {isExpanded ? (
            <TablerChevronUp className="w-5 h-5" />
          ) : (
            <TablerChevronDown className="w-5 h-5" />
          )}
        </span>
      </RACButton>
      {isExpanded && (
        <div className="px-4 py-4 border-t border-neutral-200 dark:border-neutral-800">
          <div className="space-y-2 text-text-secondary">
            {description && <p className="mb-2">{description}</p>}
            <ul className="list-disc list-inside space-y-1 ml-4">
              {links.map((link, index) => (
                <li key={index}>
                  <a
                    href={link.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-text-link underline decoration-1 underline-offset-4 hover:decoration-2 hover:underline-offset-6 transition-all"
                  >
                    {link.text}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>
      )}
    </div>
  );
}

function NumberBadge({ number }: { number: number }) {
  return (
    <div className="flex-shrink-0 w-8 h-8 rounded-full bg-button-primary text-white flex items-center justify-center font-semibold text-sm">
      {number}
    </div>
  );
}

function RouteComponent() {
  const scrollToSection = (id: string) => {
    const element = document.getElementById(id);
    element?.scrollIntoView({ behavior: "smooth", block: "start" });
  };

  return (
    <main className="px-4 sm:px-6 lg:px-8">
      <div className="prose prose-lg dark:prose-invert max-w-none">
        {/* Table of Contents */}
        <nav className="mb-12 pb-8 border-b border-neutral-200 dark:border-neutral-800">
          <Heading className="text-2xl font-bold mb-4">Resources</Heading>
          <p className="text-text-secondary mb-4">
            Find what you need based on what you’re trying to do right now.
          </p>
          <ul className="list-none space-y-2">
            {[
              ["getting-started", "Getting Started Fast"],
              ["core-features", "Building Core Features"],
              ["debugging", "Debugging & Fixing Issues"],
              ["demo-ready", "Making It Demo-Ready"],
            ].map(([id, label]) => (
              <li key={id}>
                <a
                  href={`#${id}`}
                  onClick={(e) => {
                    e.preventDefault();
                    scrollToSection(id);
                  }}
                  className="text-text-link underline decoration-1 underline-offset-4 hover:decoration-2 hover:underline-offset-8 transition-all"
                >
                  {label}
                </a>
              </li>
            ))}
          </ul>
        </nav>

        {/* Getting Started Fast */}
        <section id="getting-started" className="mb-16 scroll-mt-8">
          <Heading className="text-3xl font-bold mb-6">
            Getting Started Fast
          </Heading>
          <p className="text-text-secondary mb-6">
            Below are free and standard boilerplates that can help you get
            started with your hackathon project.
          </p>
          <div className="mb-8">
            <div className="overflow-x-auto">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="border-b border-neutral-200 dark:border-neutral-800">
                    <th className="text-left px-4 py-2 font-semibold">Name</th>
                    <th className="text-left px-4 py-2 font-semibold">
                      Technology/Framework
                    </th>
                    <th className="text-left px-4 py-2 font-semibold">Use</th>
                    <th className="text-left px-4 py-2 font-semibold">
                      Language
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {boilerplateData.map((item, index) => (
                    <tr
                      key={index}
                      className={
                        index % 2 === 0
                          ? "bg-surface"
                          : "bg-neutral-50 dark:bg-neutral-900"
                      }
                    >
                      <td className="px-4 py-2">
                        <a
                          href={item.nameLink}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-text-link underline decoration-1 underline-offset-4 hover:decoration-2 hover:underline-offset-6 transition-all"
                        >
                          {item.name}
                        </a>
                      </td>
                      <td className="px-4 py-2">
                        <a
                          href={item.technologyLink}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-text-link underline decoration-1 underline-offset-4 hover:decoration-2 hover:underline-offset-6 transition-all"
                        >
                          {item.technology}
                        </a>
                      </td>
                      <td className="px-4 py-2 text-text-secondary">
                        {item.use}
                      </td>
                      <td className="px-4 py-2 text-text-secondary">
                        {item.language}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </section>

        {/* Building Core Features */}
        <section id="core-features" className="mb-16 scroll-mt-8">
          <Heading className="text-3xl font-bold mb-6">
            Building Core Features
          </Heading>
          <p className="text-text-secondary mb-6">
            Now that you have a foundation for your project, you can start
            building core features. Here are some helpful guides for building
            core features.
          </p>

          <div className="space-y-2">
            {coreFeaturesData.map((item, index) => (
              <AccordionItem
                key={index}
                title={item.title}
                description={item.description}
                links={item.links}
              />
            ))}
          </div>
        </section>

        {/* Debugging */}
        <section id="debugging" className="mb-16 scroll-mt-8">
          <Heading className="text-3xl font-bold mb-6">
            Debugging & Fixing Issues
          </Heading>
          <p className="text-text-secondary mb-6">
            When things break or behave unexpectedly, use these resources to
            diagnose and fix issues. Don't forget to check the Discord support
            channels for real-time help from mentors and fellow hackers.
          </p>

          <div className="space-y-2">
            {debuggingData.map((item, index) => (
              <AccordionItem
                key={index}
                title={item.title}
                description={item.description}
                links={item.links}
              />
            ))}
          </div>
        </section>

        {/* Making It Demo-Ready */}
        <section id="demo-ready" className="mb-16 scroll-mt-8">
          <Heading className="text-3xl font-bold mb-6">
            Making It Demo-Ready
          </Heading>
          <p className="text-text-secondary mb-6">
            A clear demo is more memorable than a complex one.
          </p>

          <div className="space-y-6">
            <div>
              <div className="flex items-start gap-3 mb-2">
                <NumberBadge number={1} />
                <div className="flex-1">
                  <h3 className="font-semibold text-lg mb-1">
                    Start with the problem, not the tech
                  </h3>
                  <p className="text-text-secondary">
                    Clearly state who this is for and what problem you're
                    solving in a few sentences.
                  </p>
                </div>
              </div>
            </div>

            <div>
              <div className="flex items-start gap-3 mb-2">
                <NumberBadge number={2} />
                <div className="flex-1">
                  <h3 className="font-semibold text-lg mb-1">
                    Show something working early
                  </h3>
                  <p className="text-text-secondary">
                    A live UI, changing data, or physical response builds
                    credibility fast.
                  </p>
                </div>
              </div>
            </div>

            <div>
              <div className="flex items-start gap-3 mb-2">
                <NumberBadge number={3} />
                <div className="flex-1">
                  <h3 className="font-semibold text-lg mb-1">
                    Tell a short story
                  </h3>
                  <p className="text-text-secondary">
                    Problem → solution → impact. Focus on 1–2 core features, not
                    everything.
                  </p>
                </div>
              </div>
            </div>

            <div>
              <div className="flex items-start gap-3 mb-2">
                <NumberBadge number={4} />
                <div className="flex-1">
                  <h3 className="font-semibold text-lg mb-1">
                    Design for demo failure
                  </h3>
                  <p className="text-text-secondary">
                    Have screenshots, videos, or mock data ready in cases such
                    as hardware failure.
                  </p>
                </div>
              </div>
            </div>

            <div>
              <div className="flex items-start gap-3 mb-2">
                <NumberBadge number={5} />
                <div className="flex-1">
                  <h3 className="font-semibold text-lg mb-1">
                    Highlight what makes it unique
                  </h3>
                  <p className="text-text-secondary">
                    Explicitly say what's new, different, or better. Don't make
                    judges guess by making it clear and concise.
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div className="mt-8 pt-6 border-t border-neutral-200 dark:border-neutral-800">
            <h3 className="font-semibold text-lg mb-3">
              Example Presentations
            </h3>
            <ul className="space-y-2">
              <li>
                <a
                  href="https://swamphacksx.devpost.com/project-gallery"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-text-link underline decoration-1 underline-offset-4 hover:decoration-2 hover:underline-offset-8 transition-all"
                >
                  Explore SwampHacks X Winning Projects
                </a>
              </li>
              <li>
                <a
                  href="https://www.youtube.com/watch?v=5WQgLboa_I8"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-text-link underline decoration-1 underline-offset-4 hover:decoration-2 hover:underline-offset-8 transition-all"
                >
                  MIT Global AI Hackathon Presentations 2025
                </a>
              </li>
            </ul>
          </div>
        </section>

        {/* Quick References */}
        {/* <section id="quick-references" className="mb-16 scroll-mt-8">
          <Heading className="text-3xl font-bold mb-6">
            Quick References
          </Heading>
          <p className="text-text-secondary mb-6">
            Placeholder: Direct links and cheat sheets for common tasks.
          </p>

          <ul className="space-y-2">
            <li>
              <a href="#" className="text-text-link underline">
                Git & GitHub basics
              </a>
            </li>
            <li>
              <a href="#" className="text-text-link underline">
                Deployment checklist
              </a>
            </li>
            <li>
              <a href="#" className="text-text-link underline">
                Hardware pinouts
              </a>
            </li>
            <li>
              <a href="#" className="text-text-link underline">
                Command line cheatsheet
              </a>
            </li>
          </ul>
        </section> */}
      </div>
    </main>
  );
}
