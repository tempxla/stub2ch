function Logout() {
    var frm = document.getElementById("f1");
    frm.action = "/test/_admin/logout";
    frm.submit();
}

function CreateBoard(boardName) {
    var frm = document.getElementById("f1");
    frm.action = "/test/_admin/func/create-board/" + boardName;
    frm.submit();
}
