import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import userEvent from "@testing-library/user-event";
import ChooseChat from "./ChooseChat";

describe("Choose Chat", () => {
	test("render choose user chat btn", () => {
		// Arange
		render(
			<MemoryRouter>
				<ChooseChat />
			</MemoryRouter>
		);
		// Assert
		const userChatBtn = screen.getByRole("button", { name: /users/i });
		expect(userChatBtn).toBeInTheDocument();
	});
	test("render choose group chat btn", () => {
		// Arange
		render(
			<MemoryRouter>
				<ChooseChat />
			</MemoryRouter>
		);
		// Assert
		const groupChatBtn = screen.getByRole("button", { name: /groups/i });
		expect(groupChatBtn).toBeInTheDocument();
	});
	test("click choose user chat btn", () => {
		// Arange
		render(
			<MemoryRouter>
				<ChooseChat />
			</MemoryRouter>
		);
		// Act
		const userChatBtn = screen.getByRole("button", { name: /user/i });
		userEvent.click(userChatBtn);
		// Assert
		expect(userChatBtn.getAttribute("class")).toMatch("active");
	});
	test("click choose group chat btn", () => {
		// Arange
		render(
			<MemoryRouter>
				<ChooseChat />
			</MemoryRouter>
		);
		const groupChatBtn = screen.getByRole("button", { name: /groups/i });
		userEvent.click(groupChatBtn);
		// Assert
		expect(groupChatBtn.getAttribute("class")).toMatch("active");
	});
});
