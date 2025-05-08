import { forwardRef } from "react";
// RAC = React Aria Components
import { Button as RAC_Button } from "react-aria-components";

type ButtonProps = Parameters<typeof RAC_Button>[0];

const Button = forwardRef<HTMLButtonElement, ButtonProps>((props, ref) => {
  return <RAC_Button {...props} ref={ref} />;
});

export { Button };
