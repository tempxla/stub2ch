function BodyOnload(subjects, precure, boardName) {

    if (sessionStorage == undefined) {
        precure.innerHTML = "X";
        return;
    }

    var jsonText = sessionStorage.getItem(boardName + "_subject");

    if (!jsonText) {
        return;
    }
    var obj = JSON.parse(jsonText);

    var frag = createSubjectFragment(boardName, obj);
    while (subjects.firstChild) {
        subjects.removeChild(subjects.firstChild)
    }
    subjects.appendChild(frag);

    precure.innerHTML = obj.precure;
}

function GetSubject(btn, subjects, message,  last_load_time,  precure, boardName) {

    var top_load_delay     = 10 // seconds

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

    xhr.open("POST", "/" + boardName + "/subject.json");
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.send("precure=" + precure.innerHTML);

    // サーバーからクライアントへの転送の進捗 (ダウンロード)
    function updateProgress (oEvent) {
        var percentComplete;
        if (oEvent.lengthComputable) {
            percentComplete = oEvent.loaded / oEvent.total * 100;
            // ...
        } else {
            // 全体の長さが不明なため、進捗情報を計算できない
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
        var frag = createSubjectFragment(boardName, obj);
        while (subjects.firstChild) {
            subjects.removeChild(subjects.firstChild)
        }
        subjects.appendChild(frag);
        precure.innerHTML = obj.precure;

        // sessionStorage
        sessionStorage.setItem(boardName + "_subject", xhr.responseText);
    }
}

function createTdTextNode(text) {
    var td = document.createElement("td");
    var txt = document.createTextNode(text);
    td.appendChild(txt);
    return td
}

function createTdANode(link, text) {
    var td = document.createElement("td");
    var a = document.createElement("a");
    var txt = document.createTextNode(text);
    a.setAttribute("href", link);
    a.appendChild(txt);
    td.appendChild(a);
    return td
}

function createSubjectFragment(boardName, obj) {
    var frag = document.createDocumentFragment();
    for (var i = 0, len = obj.subjects.length; i < len; i++) {
        var sbj = obj.subjects[i];
        var tr = document.createElement("tr");
        var url = "/test/read.cgi/" + boardName + "/" + sbj.thread_key + "/";
        tr.appendChild(createTdANode(url, sbj.thread_title));
        tr.appendChild(createTdTextNode(sbj.message_count));
        tr.appendChild(createTdTextNode(sbj.last_modified));
        frag.appendChild(tr);
    }
    return frag;
}
