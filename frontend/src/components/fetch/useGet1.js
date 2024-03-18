import { useState, useEffect } from "react";
import axios from "axios";

// function useGet(url) {
//   const [isLoading, setIsLoading] = useState(false);
//   const [data, setData] = useState([]);
//   const [error, setError] = useState(null);

//   console.log("useGet: Starting fetch for before ue", url);
//   useEffect(() => {
//     console.log("useGet: Starting fetch for", url);
//     setIsLoading(true);
//     fetch(`http://localhost:8080${url}`)
//       .then(response => {
//         console.log("useGet: Received response:", response);
//         return response.json();
//       })
//       .then((data) => {
//         console.log("useGet: Parsed data:", data);
//         setIsLoading(false);
//         setData(data);
//       })
//       .catch((error) => {
//         console.error("useGet: Fetch Error:", error);
//         setIsLoading(false);
//         setError(error);
//       });
//   }, [url]);

//   console.log("useGet: Returning data:", data);
//   return { error, isLoading, data };
// }

// export default useGet;
// function useGet(url) {

//     const [isLoading, setIsLoading] = useState(false);
//     const [data, setData] = useState([]);
//     const [error, setError] = useState(null);

//     useEffect(() => {
//       console.log("useGet: ");
//         setIsLoading(true)
//         fetch(`http://localhost:8080${url}`)
//         .then(response => {
//           console.log("useGet: ",data);
//           return response.json()
//         })
//         .then((data) => {
//           console.log("useGet: ",data);
//             // const dataArr = [];

//             // for (const key in data) {
//             //     const value = {
//             //         id: key,
//             //         ...data[key]
//             //     };

//             //     dataArr.push(value);
//             // }
//             setIsLoading(false);
//             setData(data);
//         })
//     }, [url]);


//     return { error, isLoading, data };
// }

//Use this instead? Yes i think so

// const useGet = (url) => {
//     const [status, setStatus] = useState('idle');
//     const [data, setData] = useState([]);

//     useEffect(() => {
//         if (!url) return;
//         const fetchData = async () => {
//             setStatus('fetching');
//             const response = await fetch(`http://localhost:8080${url}`);
//             const data = await response.json();
//             setData(data);
//             setStatus('fetched');
//         };

//         fetchData();
//     }, [url]);

//     return { status, data };
// };

// export default useGet;
// /OR
// //Using axios 
const useGet = url => {
  const [data, setData] = useState([]);
  const [isLoaded, setIsLoaded] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = () => {
      axios
        .get( `http://localhost:8080${url}`, { withCredentials: true })
        .then(response => {
          setIsLoaded(true);
          setData(response.data);
        })
        .catch(error => {
          setError(error);
        });
    };
    fetchData();
  }, [url]);

  return { error, isLoaded, data };
};

export default useGet;
  
