import NextLink from "next/link";
import { MouseEventHandler, ReactNode } from "react";

export default function Link({
	children,
	href,
	className,
	onClick,
}: {
	children: ReactNode;
	href?: string;
	className?: string;
	onClick?: MouseEventHandler;
}) {
	const assignClass = `text-link hover:text-linkHover cursor-pointer ${className}`;

	if (typeof href !== "undefined") {
		return (
			<NextLink href={href} className={assignClass} onClick={onClick}>
				{children}
			</NextLink>
		);
	}
	return (
		<a className={assignClass} onClick={onClick}>
			{children}
		</a>
	);
}
