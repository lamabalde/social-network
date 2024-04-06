import { ref } from "vue";
import ws from "../ws";

export const searchComponent = ref({
  searchResults: [],
  closeFriendSearchResults: [],
  getSearchResults(query) {
    ws.request("searchGroupsUsers", query);
  },
  renderSearchResults(payload) {
    this.searchResults = payload.data;
  },
  getUserSearchResultsForGroup(query, groupID) {
    const payload = {
      user: query,
      groupID: groupID,
    };
    ws.request("searchUsersNotInGroup", payload);
  },
  renderUserSearchResults(payload) {
    this.searchResults = payload.data;
  },
  renderCloseFriendSearchResults(payload) {
    this.closeFriendSearchResults = payload.data;
  },
  searchCloseFriendsList(query) {
    ws.request("searchNotCloseFriends", query);
  },
});