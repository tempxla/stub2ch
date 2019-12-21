function logout() {
    var frm = document.getElementById("f1")
    frm.action = "/test/_admin/logout"
    frm.submit()
}
