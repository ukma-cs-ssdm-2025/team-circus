import type { SVGProps } from "react";

type DownloadIconProps = SVGProps<SVGSVGElement>;

const DownloadIcon = (props: DownloadIconProps) => (
	<svg
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 24 24"
		fill="currentColor"
		aria-hidden="true"
		{...props}
	>
		<path d="M12 16a1 1 0 0 1-.7-.29l-5-5a1 1 0 1 1 1.4-1.42l3.3 3.3V4a1 1 0 1 1 2 0v8.59l3.3-3.3a1 1 0 0 1 1.4 1.42l-5 5A1 1 0 0 1 12 16Z" />
		<path d="M5 18a1 1 0 0 1 1-1h12a1 1 0 1 1 0 2H6a1 1 0 0 1-1-1Z" />
	</svg>
);

export default DownloadIcon;
