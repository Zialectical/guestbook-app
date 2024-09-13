$(document).ready(function() {
    // Load initial guestbook entries
    loadEntries();

    // Handle form submission
    $("#guestbook-submit").click(function() {
        var entryContent = $("#guestbook-entry-content").val();
        if (entryContent) {
            submitEntry(entryContent);
            $("#guestbook-entry-content").val('');
        }
    });

    function loadEntries() {
        $.get("/lrange/guestbook", function(data) {
            $("#guestbook-entries").empty();
            data.forEach(function(entry) {
                $("#guestbook-entries").append("<p>" + entry + "</p>");
            });
        });
    }

    function submitEntry(content) {
        $.get("/rpush/guestbook/" + encodeURIComponent(content), function() {
            loadEntries();
        });
    }
});
