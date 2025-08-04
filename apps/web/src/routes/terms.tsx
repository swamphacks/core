import { createFileRoute } from "@tanstack/react-router";
import { Heading, Text } from "react-aria-components";

export const Route = createFileRoute("/terms")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="max-w-3xl mx-auto p-6">
      <Heading className="text-3xl mb-4">SwampHacks Terms of Service</Heading>
      <Text elementType="p" className="mb-6">
        Last updated August 4th, 2025
      </Text>

      <Text elementType="p" className="mb-6 text-text-main">
        Welcome to SwampHacks{" "}
        <a
          className="underline"
          href="https://github.com/SwampHacks/core"
          target="_blank"
          rel="noopener noreferrer"
        >
          Core
        </a>
        , the official hackathon platform hosted by the University of Florida.
        By registering for or using this platform, you agree to comply with
        these Terms of Service, the&nbsp;
        <a
          href="https://mlh.io/code-of-conduct"
          target="_blank"
          rel="noopener noreferrer"
          className="underline"
        >
          Major League Hacking (MLH) Code of Conduct
        </a>
        , and all applicable University of Florida policies. Please read these
        terms carefully.
      </Text>

      <Heading className="text-xl mt-6 mb-2">1. Platform Eligibility</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        Access to the SwampHacks platform is intended for students and
        organizers affiliated with hackathons hosted or supported by the
        University of Florida. By using this platform, you affirm that you meet
        any eligibility criteria set forth by the individual events you
        participate in.
      </Text>

      <Heading className="text-xl mt-6 mb-2">2. Code of Conduct</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        SwampHacks promotes a respectful, inclusive, and harassment-free
        environment in line with the MLH Code of Conduct. Users of this platform
        must behave accordingly. Any violation may result in suspension or
        removal from the platform or associated events without refund or appeal.
      </Text>

      <Heading className="text-xl mt-6 mb-2">3. Intellectual Property</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        Users retain ownership of all intellectual property submitted or shared
        on the platform. By uploading content, you grant SwampHacks and event
        organizers a non-exclusive, worldwide license to use and display your
        submissions for promotional, educational, and administrative purposes.
      </Text>

      <Heading className="text-xl mt-6 mb-2">4. Privacy and Data Usage</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        SwampHacks collects and processes personal information necessary to
        facilitate hackathon participation and platform functionality, following
        University of Florida policies and applicable laws. Please review
        our&nbsp;
        <a href="/privacy" className="underline">
          Privacy Policy
        </a>
        &nbsp;for detailed information.
      </Text>

      <Heading className="text-xl mt-6 mb-2">
        5. Platform Administration
      </Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        SwampHacks administrators reserve the right to manage platform content,
        enforce rules, and suspend or remove users or content that violate these
        Terms or other applicable policies. Final decisions regarding platform
        use, content, and disputes rest with SwampHacks and its administrators.
      </Text>

      <Heading className="text-xl mt-6 mb-2">6. Event-Specific Rules</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        While this platform provides tools for hackathon registration and
        project submission, each individual event hosted on SwampHacks may have
        its own specific rules, judging criteria, and policies. Users are
        responsible for reviewing and adhering to the rules set by each event
        organizer.
      </Text>

      <Heading className="text-xl mt-6 mb-2">
        7. Liability and Disclaimer
      </Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        SwampHacks, its platform administrators, sponsors, and the University of
        Florida disclaim all liability for any loss, injury, or damages arising
        from use of the platform or participation in events. Users assume all
        risks related to their platform usage and event participation.
      </Text>

      <Heading className="text-xl mt-6 mb-2">8. Changes to Terms</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        SwampHacks reserves the right to modify these Terms of Service at any
        time. Users will be notified of material changes. Continued use of the
        platform after updates constitutes acceptance of the revised terms.
      </Text>

      <Text elementType="p" className="mt-10 underline">
        By using the SwampHacks Core platform, you acknowledge that you have
        read, understood, and agreed to these Terms of Service.
      </Text>
    </div>
  );
}
