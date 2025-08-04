import { createFileRoute } from "@tanstack/react-router";
import { Heading, Text } from "react-aria-components";

export const Route = createFileRoute("/privacy")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="max-w-3xl mx-auto p-6">
      <Heading className="text-3xl mb-4">SwampHacks Privacy Policy</Heading>
      <Text elementType="p" className="mb-6">
        Last updated August 4th, 2025
      </Text>

      <Text elementType="p" className="mb-6 text-text-main">
        SwampHacks (“we”, “our”, or “us”) values your privacy. This Privacy
        Policy explains how we collect, use, and protect your personal
        information when you use the SwampHacks platform. By registering for or
        participating in SwampHacks events or using our platform, you consent to
        the data practices described in this policy.
      </Text>

      <Heading className="text-xl mt-6 mb-2">1. Information We Collect</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        We collect personal information you provide during registration and
        platform use, including but not limited to:
      </Text>
      <ul className="list-disc list-inside mb-4 text-text-secondary">
        <li>Full name, email address, and Discord profiles</li>
        <li>Educational affiliation and enrollment status</li>
        <li>
          Social media profiles or handles you provide for identification or
          outreach
        </li>
        <li>Project submissions and related content</li>
      </ul>

      <Heading className="text-xl mt-6 mb-2">
        2. How We Use Your Information
      </Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        Your information is used to:
      </Text>
      <ul className="list-disc list-inside mb-4 text-text-secondary">
        <li>Verify participant eligibility and identity</li>
        <li>
          Communicate important event updates, announcements, and outreach
        </li>
        <li>Facilitate project judging and event administration</li>
        <li>Comply with legal and University of Florida policies</li>
      </ul>

      <Heading className="text-xl mt-6 mb-2">
        3. Data Sharing and Disclosure
      </Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        We will never sell your personal data. However, we may share your
        information with:
      </Text>
      <ul className="list-disc list-inside mb-4 text-text-secondary">
        <li>
          Event judges and organizers for project evaluation and communication
        </li>
        <li>
          Major League Hacking (MLH) as part of partnership and event
          coordination
        </li>
        <li>
          Other trusted third parties who assist with platform operations or
          compliance
        </li>
      </ul>
      <Text elementType="p" className="mb-4 text-text-secondary">
        All third parties are required to handle your data securely and only use
        it for authorized purposes.
      </Text>

      <Heading className="text-xl mt-6 mb-2">4. Data Security</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        We implement reasonable administrative, technical, and physical
        safeguards to protect your personal information against unauthorized
        access, disclosure, alteration, and destruction.
      </Text>

      <Heading className="text-xl mt-6 mb-2">5. Your Rights</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        Depending on your location, you may have certain rights regarding your
        personal data, including the right to access, correct, or delete your
        information. To exercise these rights or for any privacy-related
        inquiries, please contact us at{" "}
        <a href="mailto:contact@swamphacks.com" className="underline">
          contact@swamphacks.com
        </a>
        .
      </Text>

      <Heading className="text-xl mt-6 mb-2">6. Changes to This Policy</Heading>
      <Text elementType="p" className="mb-4 text-text-secondary">
        We may update this Privacy Policy periodically. We will notify users of
        significant changes via email or platform announcements. Continued use
        of SwampHacks after changes constitutes acceptance of the updated
        policy.
      </Text>

      <Text elementType="p" className="mt-10 underline">
        By using SwampHacks, you acknowledge that you have read and agree to
        this Privacy Policy.
      </Text>
    </div>
  );
}
