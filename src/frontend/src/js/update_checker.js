const events = ["startupUpdateCheckerError", "startupUpdateCheckerUpdateAvailable"];

const unlistenEvents = () => window.runtime.EventsOff(events[0], events[1]);
const eventsEnd = () => {
    window.runtime.EventsEmit("ready");
};

window.runtime.EventsOnce(events[0], async () => {
    unlistenEvents();

    // TODO: ask user whether to continue or not

    eventsEnd();
});
window.runtime.EventsOnce(events[1], async (latestVersion) => {
    unlistenEvents();

    console.log(latestVersion);

    // TODO: ask user, if accepted, redirect

    eventsEnd();
});