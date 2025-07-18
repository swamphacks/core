/* eslint-disable */
import { type ReactNode } from "react";
import { TextField } from "@/components/ui/TextField";
import { Form } from "react-aria-components";
import { Button } from "@/components/ui/Button";
import { FieldGroup } from "@/components/ui/Field";
import TablerPhone from "~icons/tabler/phone";
import { Radio, RadioGroup } from "@/components/ui/Radio";
import { Select, SelectItem } from "@/components/ui/Select";
import TablerAt from "~icons/tabler/at";
import TablerBrandLinkedin from "~icons/tabler/brand-linkedin";
import TablerBrandGithub from "~icons/tabler/brand-github";
import { ComboBox, ComboBoxItem } from "@/components/ui/ComboBox";
import { MultiSelect } from "@/components/ui/MultiSelect";

interface ApplicationFormProps {
  title: ReactNode;
  description?: ReactNode;
}

const ApplicationForm = ({ title, description }: ApplicationFormProps) => {
  let onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    // Prevent default browser page refresh.
    e.preventDefault();

    // Get form data as an object.
    let data = Object.fromEntries(new FormData(e.currentTarget));
    console.log(data);
  };

  return (
    <div className="w-full sm:max-w-180 mx-auto font-figtree pb-2 p-2">
      <div className="space-y-3 py-5 border-b-1 border-border">
        <p className="text-2xl text-text-main">{title}</p>
        <p className="text-text-secondary">{description}</p>
      </div>

      <Form className="space-y-3 mt-5" onSubmit={onSubmit}>
        <div className="pb-5">
          <p className="text-lg">Personal Information</p>

          <div className="mt-3 space-y-4">
            <FieldGroup className="gap-4">
              <TextField
                label="First Name"
                placeholder="Enter your first name"
                isRequired
                className="flex-1"
                name="firstName"
              />
              <TextField
                label="Last Name"
                placeholder="Enter your last name"
                isRequired
                className="flex-1"
                name="lastName"
              />
            </FieldGroup>
            {/* <FieldGroup className="gap-4">
              <MultiSelect name="test" isRequired isMulti />
              <TextField
                label="Age"
                placeholder="Enter your age"
                isRequired
                className="flex-1"
                name="age"
              />
              <Select
                label="Level of Study"
                placeholder="Select your level"
                isRequired
                className="flex-1"
                name="testaaa"
              >
                <SelectItem>1st Year</SelectItem>
                <SelectItem>2nd Year</SelectItem>
                <SelectItem>3rd Year</SelectItem>
                <SelectItem>4th Year</SelectItem>
                <SelectItem>Grad Student</SelectItem>
              </Select>
            </FieldGroup> */}

            <FieldGroup className="gap-4">
              <TextField
                label="Age"
                placeholder="Enter your age"
                isRequired
                className="flex-1"
              />
              <TextField
                label="Phone"
                placeholder="Enter a phone number"
                isRequired
                className="flex-1"
                icon={TablerPhone}
              />
            </FieldGroup>

            <FieldGroup className="gap-4">
              <TextField
                label="Preferred Email"
                placeholder="Enter your email"
                isRequired
                className="flex-1"
                icon={TablerAt}
              />
              <TextField
                label="University Email"
                placeholder="Enter your school email"
                isRequired
                className="flex-1"
                icon={TablerAt}
              />
            </FieldGroup>

            <FieldGroup className="gap-4">
              <TextField
                label="Linkedin URL"
                placeholder="Enter link"
                isRequired
                className="flex-1"
                icon={TablerBrandLinkedin}
              />
            </FieldGroup>

            <FieldGroup className="gap-4">
              <TextField
                label="Github URL"
                placeholder="Enter link"
                isRequired
                className="flex-1"
                icon={TablerBrandGithub}
              />
            </FieldGroup>
          </div>
        </div>

        <div className="pb-5">
          <p className="text-lg">Education</p>
          <div className="mt-3 space-y-4">
            <FieldGroup className="gap-4">
              <ComboBox
                label="School"
                placeholder="Select your school"
                isRequired
                className="flex-1"
              >
                <ComboBoxItem>Chocolate</ComboBoxItem>
                <ComboBoxItem>Mint</ComboBoxItem>
                <ComboBoxItem>Strawberry</ComboBoxItem>
                <ComboBoxItem>Vanilla</ComboBoxItem>
              </ComboBox>
              <Select
                label="Level of Study"
                placeholder="Select your level"
                isRequired
                className="flex-1"
              >
                <SelectItem>1st Year</SelectItem>
                <SelectItem>2nd Year</SelectItem>
                <SelectItem>3rd Year</SelectItem>
                <SelectItem>4th Year</SelectItem>
                <SelectItem>Grad Student</SelectItem>
              </Select>
            </FieldGroup>

            <FieldGroup className="gap-4">
              <TextField
                label="Year in College"
                placeholder="Select your year"
                isRequired
                className="flex-1"
              />
              <TextField
                label="Graduation Date"
                placeholder="Select date"
                isRequired
                className="flex-1"
              />
            </FieldGroup>

            <FieldGroup className="gap-4">
              {/* @ts-ignore */}
              <MultiSelect name="Major(s)" isRequired isMulti />
            </FieldGroup>

            <FieldGroup className="gap-4">
              {/* @ts-ignore */}
              <MultiSelect name="Minors(s)" isRequired isMulti />
            </FieldGroup>
          </div>
        </div>

        <div>
          <p className="text-lg">Experience & Preferences</p>
          <div className="mt-3 space-y-4">
            <FieldGroup className="gap-4">
              <RadioGroup
                label="How many hackathons have you participated in?"
                isRequired
                className="grow"
              >
                <Radio value="0">SwampHacks would be my first</Radio>
                <Radio value="1">1</Radio>
                <Radio value="2">2</Radio>
                <Radio value="3">3</Radio>
                <Radio value="4">4</Radio>
                <Radio value="5">5</Radio>
                <Radio value="6">6</Radio>
              </RadioGroup>
            </FieldGroup>

            <FieldGroup className="gap-4">
              <Select
                label="T-Shirt Size"
                placeholder="Select your size"
                isRequired
                className="flex-1"
              >
                <SelectItem>Small</SelectItem>
                <SelectItem>Medium</SelectItem>
                <SelectItem>Large</SelectItem>
              </Select>
            </FieldGroup>
          </div>
        </div>

        <div className="flex justify-end">
          <Button type="submit">Submit</Button>
        </div>
      </Form>
    </div>
  );
};

export { ApplicationForm };
