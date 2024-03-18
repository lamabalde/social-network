import { useState, useEffect } from "react";

function useGet(url) {
  const [isLoading, setIsLoading] = useState(false);
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);

  console.log("useGet: Starting fetch for before ue", url);
  useEffect(() => {
    console.log("useGet: Starting fetch for", url);
    setIsLoading(true);
    fetch(`http://localhost:8080${url}`, {credentials: "include"})
      .then(response => {
        console.log("useGet: Received response:", response);
        return response.json();
      })
      .then((data) => {
        console.log("useGet: Parsed data:", data);
        setIsLoading(false);
        setData(data);
      })
      .catch((error) => {
        console.error("useGet: Fetch Error:", error);
        setIsLoading(false);
        setError(error);
      });
  }, [url]);

  console.log("useGet: Returning data:", data);
  return { error, isLoading: isLoading, data };
}

export default useGet;