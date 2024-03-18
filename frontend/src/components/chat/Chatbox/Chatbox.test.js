import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import Chatbox from "./Chatbox";

describe("Old Chat Msg", () => {
	test("Loading Old Msgs (New Chat)", async () => {
		// Assign
		window.fetch = jest.fn();
		window.fetch.mockResolvedValueOnce({
			json: async () => [],
		});
		// Act
		//Assert
		const newChatMsg = await screen.findByText(
			/Start Chatting/i,
			{ exact: false },
			{ timeout: 3000 }
		);
		expect(newChatMsg).toBeInTheDocument();
	});

	test("Loading Old Msgs (Old Msg >= 1 )", async () => {
		// Assign
		window.fetch = jest.fn();
		window.fetch.mockResolvedValueOnce({
			json: async () => {
				{
					data: [
						{
							createdat: "2023-03-16 00:00:00.0000000 +0100 +0100",
							id: 1,
							label: "",
							message: "testOldMsg0",
							onlineuserids: null,
							sourceid: 1,
							targetid: 2,
						},
						{
							createdat: "2023-03-16 00:01:00.0000000 +0100 +0100",
							id: 2,
							label: "",
							message: "testOldMsg1",
							onlineuserids: null,
							sourceid: 2,
							targetid: 1,
						},
					];
				}
			},
		});
		// Act
		//Assert
		const msg0 = await screen.findByText("testOldMsg0", {}, { timeout: 3000 });
		expect(msg0).toBeInTheDocument();
		const msg1 = await screen.findByText("testOldMsg1", {}, { timeout: 3000 });
		expect(msg1).toBeInTheDocument();
	});
});
