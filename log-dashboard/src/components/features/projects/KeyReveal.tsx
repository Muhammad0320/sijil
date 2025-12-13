interface IKeyRevel {
  data: {
    projectId: number;
    name: string;
    apiKey: string;
    apiSecret: string;
  };
}

export default function KeyRevel({
  data: { projectId, name, apiKey, apiSecret },
}: IKeyRevel) {
  return (
    <div>
      <h3> Keys for project: {name} </h3>
      <p>
        {" "}
        Project id: <span>{projectId}</span>{" "}
      </p>{" "}
      {/* Why do we even need to pass the project id not to talk of revealing it to the user? Should it be removed? */}
      <p>
        {" "}
        API key:
        <span> {apiKey} </span>
        <button> copy icon </button>
      </p>
      <p>
        {" "}
        API secret:
        <span> {apiSecret} </span>
        <button> copy icon </button>
      </p>
    </div>
  );
}
