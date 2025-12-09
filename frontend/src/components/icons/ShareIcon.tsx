type ShareIconProps = {
	className?: string;
};

const ShareIcon = ({ className }: ShareIconProps) => (
	<svg
		className={className}
		viewBox="0 0 24 24"
		fill="none"
		xmlns="http://www.w3.org/2000/svg"
		aria-hidden
	>
		<path
			d="M8.5 11.5L14.5 8M8.5 12.5L14.5 16"
			stroke="currentColor"
			strokeWidth="1.6"
			strokeLinecap="round"
			strokeLinejoin="round"
		/>
		<circle
			cx="6.5"
			cy="12"
			r="2.25"
			stroke="currentColor"
			strokeWidth="1.6"
		/>
		<circle
			cx="17.5"
			cy="7"
			r="2.25"
			stroke="currentColor"
			strokeWidth="1.6"
		/>
		<circle
			cx="17.5"
			cy="17"
			r="2.25"
			stroke="currentColor"
			strokeWidth="1.6"
		/>
	</svg>
);

export default ShareIcon;
