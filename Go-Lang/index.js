const { GraphQLServer } = require("graphql-yoga");
const fetch = require("node-fetch");

const typeDefs = `
type Query {
        hello(name: String): String!
        getService(id: Int!): Service
    }

    type Service {
        basePrice: Int,
        cost: Int,
        cpu: Int,
        id: Int,
        name: String,
        quantity: Int,
        ram: Int
    }
`;

const resolvers = {
  Query: {
    hello: (_, { name }) => `Hello ${name || "World"}`,
    getService: async (_, { id }) => {
        const response = await fetch(`http://localhost:8081/cost/${id}`);
        return response.json();
    }
    // ,
    // getServices: async (_,{}) => {
    //     const response = await fetch(`http://localhost:8081/cost-add`);
    //     return response.json();
    // } 
    }
};

const server = new GraphQLServer({ typeDefs, resolvers });
server.start(() => console.log("Server is running on localhost:4000"));