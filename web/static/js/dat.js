function BodyOnload(boardName, threadKey) {

    var title = document.getElementById('title');
    var messages = document.getElementById('messages');
    var precure = document.getElementById('precure');

    if (localStorage == undefined) {
        precure.innerHTML = "X";
        return;
    }

    // display thread title
    var jsonText = localStorage.getItem(boardName + "_subject");
    if (jsonText != null) {
        var obj = JSON.parse(jsonText);
        for (var i = 0, len = obj.subjects.length; i < len; i++) {
            if (obj.subjects[i].thread_key == threadKey) {
                document.title = obj.subjects[i].thread_title;
                title.innerHTML = obj.subjects[i].thread_title;
                break;
            }
        }
    }

    // display cache
    var jsonDat = localStorage.getItem(boardName + "_" + threadKey);
    if (!jsonDat) {
        return;
    }
    var dat = JSON.parse(jsonDat);

    var frag = createMessageFragment(dat);

    // fix me >>>
    while (messages.firstChild) {
        messages.removeChild(messages.firstChild);
    }
    messages.appendChild(frag);
    // fix me <<<

    messages.appendChild(frag);

    document.title = dat.thread_title;
    title.innerHTML = dat.thread_title;

    // which precure
    var p = localStorage.getItem("precure");
    if (p != null) {
        precure.innerHTML = p;
    }
}

function GetDat(btn, boardName, threadKey) {

    var messages = document.getElementById('messages');
    var message = document.getElementById('message');
    var last_load_time = document.getElementById('last_load_time');
    var precure = document.getElementById('precure');

    var top_load_delay     = 10; // seconds

    if (btn.innerHTML != "LOAD") {
        return;
    }
    var timer = function(count) {
        return function() {
            if (count > 0) {
                btn.innerHTML = "LOAD   (wait " + count + " ...)";
                setTimeout(timer(count-1), 1000);
            } else {
                btn.innerHTML = "LOAD";
            }
        }
    };
    setTimeout(timer(top_load_delay), 0);

    message.innerHTML = "";
    last_load_time.innerHTML = "";

    var xhr = new XMLHttpRequest();

    xhr.addEventListener("progress", updateProgress);
    xhr.addEventListener("load", transferComplete);

    xhr.open("POST", "/" + boardName + "/json/" + threadKey + ".json");
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

    // 鯖に優しく
    var storePrecure = localStorage.getItem("precure");
    if (storePrecure != null) {
        xhr.send("precure=" + storePrecure);
    } else {
        xhr.send("precure=" + precure.innerHTML);
    }

    // Fowarding progress from server to client. (downloading)
    function updateProgress (oEvent) {
        var percentComplete;
        if (oEvent.lengthComputable) {
            percentComplete = oEvent.loaded / oEvent.total * 100;
            // ...
        } else {
            // It cannot calc progress because of unknown total length.
            percentComplete = -1;
        }
        message.innerHTML = percentComplete + " % load done.";
    }

    function transferComplete(evt) {

        message.innerHTML = "100 % load done.";
        last_load_time.innerHTML = "Last Loaded: " + new Date();
        if (xhr.status != 200) {
            message.innerHTML = xhr.status + ": " + xhr.responseText;
            return;
        }

        var obj = JSON.parse(xhr.responseText);
        var frag = createMessageFragment(obj);

        // fix me >>>
        while (messages.firstChild) {
            messages.removeChild(messages.firstChild);
        }
        messages.appendChild(frag);
        // fix me <<<

        messages.appendChild(frag);

        document.title = obj.thread_title;
        title.innerHTML = obj.thread_title;

        precure.innerHTML = obj.precure;

        // localStorage
        localStorage.setItem(boardName + "_" + threadKey, xhr.responseText);
        localStorage.setItem("precure", obj.precure);
    }
}

function createDivTextNode(clazz, text) {
    var div = document.createElement("div");
    div.setAttribute("class", clazz);
    var txt = document.createTextNode(text);
    div.appendChild(txt);
    return div;
}

function createDivNode(clazz) {
    var div = document.createElement("div");
    div.setAttribute("class", clazz);
    return div;
}

function createMessageFragment(obj) {

    var display_count = 10;

    var frag = document.createDocumentFragment();
    for (var i = 0, len = obj.messages.length; i < len && i < display_count; i++) {
        var msg = obj.messages[i];
        // name, date, id
        var row = createDivNode("row");
        row.appendChild(createDivInnerHtml("eight columns",
                                          msg.num + ": " +
                                          "<b>" + msg.name + "</b> " +
                                          "[" + msg.mail + "]"));
        row.appendChild(createDivTextNode("four columns", msg.date_and_id));
        frag.appendChild(row);
        // message
        row = createDivInnerHtml("row", msg.content);
        frag.appendChild(row);
        // br
        frag.appendChild(document.createElement("br"));
    }

    // sorry
    if (obj.messages.length > display_count) {
        var div = document.createElement("div");
        div.innerHTML = "<b>" + display_count + "までしか見せられないよ＞＜;;" + "</b>";
        frag.appendChild(div);
        frag.appendChild(document.createElement("br"));
    }

    return frag;
}

function createDivInnerHtml(clazz, text) {
    var div = document.createElement("div");
    div.setAttribute("class", clazz);
    div.innerHTML = text;
    return div;
}
