window.onbeforeunload = function(event) {

    if (hasData()) {
        return "Warning message";//warning message
    }
};

function hasData() {
    let someThing = document.getElementById("someThing");

    return someThing.value !== "submitted";


}