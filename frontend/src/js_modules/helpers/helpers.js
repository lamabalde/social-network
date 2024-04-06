import { postsComponent } from "../reactive_elements/postComponent";

export async function formToObj(ev) {
    const formObj = {}
    const formData = new FormData(ev.target)

    try {
      for (const [key, formEntry] of formData) {
        // using for loop instead of foreach for await functionality
        if (key == "password") {
          var hashedPassword = md5(formEntry);
          formObj[key] = hashedPassword; // byte array of hashed password
        } else if (key == "image") {
          if (!formEntry.name == "") {
            const imgData = formEntry.type.split("/");
            if (imgData.includes("image")) {
              formObj["imageType"] = imgData[1]; // for displaying img right after posting it
              const encodedImg = await encodeImageFileAsURL(formEntry);
              formObj[key] = encodedImg;
            } else {
              throw new Error("invalid file type")
            }
          }
        } else {
          formObj[key] = formEntry;
        }
      }
      return formObj;
    } catch (err) {
      return err
    }
}

export function convertDateToLocale(date) {
    const convertedDate = new Date(date)
    return convertedDate.toLocaleString("en-GB", {
      timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone, // gets the user's timezone
    });
}

export function convertDateToDateString(date) {
  const convertedDate = new Date(date);
  return convertedDate.toDateString();
}

export function encodeImageFileAsURL(file) {
  var reader = new FileReader();
  return new Promise(resolve => {
    reader.onloadend = function () {
        resolve(reader.result)
    };
    reader.readAsDataURL(file);
  })
}

export function findIndexInArray(array, item, argumentToCompare) {
  const index = array.findIndex((i) => {
    if (argumentToCompare!= undefined) {
      return i[argumentToCompare] === item;
    } else {
      return i === item
    }
  });
  return index
}

export function sortObjArrayByParameter(array, parameters) {
  array.sort(function (a, b) {
    return a[parameters] > b[parameters];
  });
  return array
}