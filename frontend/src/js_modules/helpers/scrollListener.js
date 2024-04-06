import { ref } from "vue";

export const atBottom = ref(false)
export const atTop = ref(false);

export function scrollAtBottomListener() {
    document.addEventListener("scroll", function (e) {
    let documentHeight = document.body.scrollHeight;
    let currentScroll = window.scrollY + window.innerHeight;
    // When the user is [modifier]px from the bottom, fire the event.
    let modifier = 200;
    if (currentScroll + modifier > documentHeight && !atBottom.value) {
        atBottom.value = true;
    } else if (!(currentScroll + modifier > documentHeight)) {
        atBottom.value = false;
    }
    });
}

